package auth

import (
	"context"
	"errors"
	"ginson/api/response"
	"ginson/config"
	"ginson/core/const/cache"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const JWTAuthHeader = "Authorization"

var (
	TokenExpired     = errors.New("token expired")
	TokenNotValidYet = errors.New("token not valid yet")
	TokenMalformed   = errors.New("token malformed")
	TokenInvalid     = errors.New("token invalid")
)

func JWT(ctx *gin.Context) {
	token, _ := ctx.Cookie(JWTAuthHeader)
	if token == "" {
		token = ctx.GetHeader(JWTAuthHeader)
	}
	if token == "" {
		response.Unauthorized(ctx, errors.New("unauthorized"))
		return
	}
	j := NewJWTAuth()
	j.Token = strings.TrimPrefix(token, "Bearer ")
	err := j.ParseToken()
	if err != nil {
		response.Unauthorized(ctx, err.Error())
		return
	}
	ctx.Set("claims", j.Claims)
	rt, _ := parseDuration(config.Config.Auth.JWT.RefreshTime)
	if j.Claims.ExpiresAt.Unix()-time.Now().Unix() < int64(rt.Seconds()) {
		log.Debug().Msgf("refresh token")
		j.RefreshToken()
		j.SetCookie(ctx)
	}
	ctx.Next()
}

type CustomClaims struct {
	BaseClaims
	jwt.RegisteredClaims
}

type BaseClaims struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username    string             `json:"username" bson:"username"`
	Nickname    string             `json:"nickname" bson:"nickname"`
	AuthorityId uint               `json:"authorityId" bson:"authorityId"`
}

type JWTAuth struct {
	SigningKey []byte
	Claims     *CustomClaims
	Token      string
}

func NewJWTAuth() *JWTAuth {
	return &JWTAuth{
		SigningKey: []byte(config.Config.Auth.JWT.SigningKey),
	}
}

func (j *JWTAuth) CreateClaims(baseClaims BaseClaims) {
	et, _ := parseDuration(config.Config.Auth.JWT.ExpiresTime)
	j.Claims = &CustomClaims{
		BaseClaims: baseClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"GINSON"},
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(et)),
			Issuer:    "GINSON",
		},
	}
	return
}

func (j *JWTAuth) GetToken() (err error) {
	if cache.Redis != nil {
		j.Token, err = j.GetCache()
		if err != nil {
			log.Error().Msgf("get cache failed: %v", err)
			err = errors.New("get cache failed")
			return
		}
		if j.Token != "" {
			return
		}
	}
	err = j.CreateToken()
	if err != nil {
		log.Error().Msgf("create token failed: %s", err.Error())
		err = errors.New("create token failed")
		return
	}
	if cache.Redis != nil {
		err = j.SetCache()
		if err != nil {
			log.Error().Msgf("set cache failed: %v", err)
			err = errors.New("set cache failed")
			return
		}
	}
	return
}

func (j *JWTAuth) CreateToken() (err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j.Claims)
	j.Token, err = token.SignedString(j.SigningKey)
	return
}

func (j *JWTAuth) ParseToken() (err error) {
	token, err := jwt.ParseWithClaims(j.Token, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return TokenNotValidYet
			} else {
				return TokenInvalid
			}
		}
	}
	if token != nil && token.Valid {
		j.Claims = token.Claims.(*CustomClaims)
		return
	} else {
		return TokenInvalid
	}
}

func (j *JWTAuth) RefreshToken() {
	et, _ := parseDuration(config.Config.Auth.JWT.ExpiresTime)
	j.Claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(et))
	j.Claims.RegisteredClaims.NotBefore = jwt.NewNumericDate(time.Now().Add(-1000))
	err := j.CreateToken()
	if err != nil {
		log.Error().Msgf("create token failed: %s", err.Error())
		err = errors.New("create token failed")
		return
	}
	if cache.Redis != nil {
		err = j.SetCache()
		if err != nil {
			log.Error().Msgf("set cache failed: %v", err)
			err = errors.New("set cache failed")
			return
		}
	}
}

func (j *JWTAuth) SetCache() (err error) {
	et, err := parseDuration(config.Config.Auth.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	err = cache.Redis.Set(context.Background(), j.Claims.Username, j.Token, et).Err()
	return
}

func (j *JWTAuth) GetCache() (token string, err error) {
	token, err = cache.Redis.Get(context.Background(), j.Claims.Username).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
			return
		}
	}
	return
}

func (j *JWTAuth) SetCookie(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     JWTAuthHeader,
		Value:    "Bearer " + j.Token,
		MaxAge:   int(j.Claims.RegisteredClaims.ExpiresAt.Unix() - time.Now().Unix()),
		Path:     "/",
		HttpOnly: false,
		Secure:   false,
	})
}

func parseDuration(d string) (time.Duration, error) {
	d = strings.TrimSpace(d)
	dr, err := time.ParseDuration(d)
	if err == nil {
		return dr, nil
	}
	if strings.Contains(d, "d") {
		index := strings.Index(d, "d")
		hour, _ := strconv.Atoi(d[:index])
		dr = time.Hour * 24 * time.Duration(hour)
		ndr, err := time.ParseDuration(d[index+1:])
		if err != nil {
			return dr, nil
		}
		return dr + ndr, nil
	}
	dv, err := strconv.ParseInt(d, 10, 64)
	return time.Duration(dv), err
}
