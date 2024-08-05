package service

import (
	"errors"
	"ginson/api/request"
	"ginson/api/response"
	"ginson/middleware"
	"ginson/middleware/auth"
	"ginson/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var User = userService{}

type userService struct {
}

func (us *userService) Login(ctx *gin.Context, login *request.Login) (err error) {
	var user model.User
	err = user.DB().FindOne(ctx, bson.M{"username": login.Username}).Decode(&user)
	if err != nil {
		err = errors.New("user not found")
		return
	}
	log.Debug().Msgf("get user info: %v", user)
	if user.Locked {
		err = errors.New("user is locked")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		err = errors.New("wrong password")
		return
	}
	jwtAuth := auth.NewJWTAuth()
	jwtAuth.CreateClaims(auth.BaseClaims{
		ID:          user.Id,
		Username:    user.Username,
		Nickname:    user.Nickname,
		AuthorityId: 0,
	})
	err = jwtAuth.GetToken()
	jwtAuth.SetCookie(ctx)
	return
}

func (us *userService) Create(ctx *gin.Context, user *model.User) (id primitive.ObjectID, err error) {
	cursor, err := user.DB().Find(ctx, bson.M{"username": user.Username})
	if err != nil {
		return
	}
	if cursor.RemainingBatchLength() > 0 {
		return id, errors.New("user already exists")
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(password)
	log.Info().Msgf("user create: %v\n", user.Username)
	user.Defaults()
	_, err = user.DB().InsertOne(ctx, user)
	return user.Id, err
}

func (us *userService) Patch(ctx *gin.Context, id primitive.ObjectID, user *model.User) (err error) {
	user.DefaultUpdateAt()
	_, err = user.DB().UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": user})
	return
}

func (us *userService) Get(ctx *gin.Context, id primitive.ObjectID) (user *model.User, err error) {
	user = &model.User{}
	err = user.DB().FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return
}

func (us *userService) Delete(ctx *gin.Context, id primitive.ObjectID) (err error) {
	user := &model.User{}
	_, err = user.DB().DeleteOne(ctx, bson.M{"_id": id})
	return
}

func (us *userService) List(ctx *gin.Context, filter bson.M, opts ...*options.FindOptions) (list *response.List, err error) {
	user := &model.User{}
	total, err := user.DB().CountDocuments(ctx, filter)
	if err != nil {
		return
	}
	cursor, err := user.DB().Find(ctx, filter, opts...)
	if err != nil {
		return
	}
	items := make([]model.User, 0)
	err = cursor.All(ctx, &items)
	if err != nil {
		return
	}
	page := middleware.GetPageFromOptions(opts...)
	list = &response.List{
		Items:   items,
		Total:   total,
		Page:    page.Page,
		PerPage: page.PerPage,
	}
	return
}
