package api

import (
	"errors"
	"ginson/api/request"
	"ginson/api/response"
	"ginson/middleware"
	"ginson/middleware/auth"
	"ginson/model"
	"ginson/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	apis = append(apis, &User{})
}

type User struct {
}

func (u *User) Prefix() string {
	return "/user"
}

func (u *User) Router(router *gin.RouterGroup) {
	router.POST("login", u.Login)
	router.POST("", u.Post)
	router.Use(auth.JWT)
	router.Group("").Use(middleware.Paginator, middleware.Sorter, middleware.Filter).GET("", u.List)
	router.GET(":id", u.Get)
	router.PATCH(":id", u.Patch)
	router.DELETE(":id", u.Delete)
}

// Login godoc
//
// @Summary	Login user with username and password
// @Tags		User
// @Param		login			body		request.Login				true		"Login Info"
// @Success		200				{object}	response.Response
// @Failure		400				{object}	response.Response
// @Failure		401				{object}	response.Response
// @Failure		500				{object}	response.Response
// @Router	/user/login [post]
func (u *User) Login(ctx *gin.Context) {
	var login request.Login
	err := ctx.ShouldBindJSON(&login)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	err = validate.Struct(login)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	err = service.User.Login(ctx, &login)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Accepted(ctx)
}

// List godoc
//
// @Summary	List user by id
// @Tags		User
// @Success		200				{object}	response.Response
// @Failure		400				{object}	response.Response
// @Failure		401				{object}	response.Response
// @Failure		500				{object}	response.Response
// @Router	/user [get]
func (u *User) List(ctx *gin.Context) {
	list, err := service.User.List(ctx, middleware.GetFilterOpt(ctx), middleware.GetSortOpt(ctx), middleware.GetPageOpt(ctx))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Data(ctx, list)
}

// Get godoc
//
// @Summary	Get user by id
// @Tags		User
// @Param		id				path		string				true		"User id"
// @Success		200				{object}	response.Response
// @Failure		400				{object}	response.Response
// @Failure		401				{object}	response.Response
// @Failure		500				{object}	response.Response
// @Router	/user/{id} [get]
func (u *User) Get(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, errors.New("id is not valid"))
		return
	}
	user, err := service.User.Get(ctx, id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Data(ctx, user)
}

// Post godoc
//
// @Summary	Post user by id
// @Tags		User
// @Param		user			body		request.User			true		"User Info"
// @Success		200				{object}	response.Response
// @Failure		400				{object}	response.Response
// @Failure		401				{object}	response.Response
// @Failure		500				{object}	response.Response
// @Router	/user [post]
func (u *User) Post(ctx *gin.Context) {
	var user request.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	_, err = service.User.Create(ctx, &model.User{
		Username: user.Username,
		Password: user.Password,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Phone:    user.Phone,
		Email:    user.Email,
	})
	if err != nil {
		response.BadRequest(ctx, err.Error())
	} else {
		response.Created(ctx)
	}
}

// Patch godoc
//
// @Summary	Patch user by id
// @Tags		User
// @Param		id				path		string				true		"User id"
// @Param		user			body		request.User			true		"User Info"
// @Success		200				{object}	response.Response
// @Failure		400				{object}	response.Response
// @Failure		401				{object}	response.Response
// @Failure		500				{object}	response.Response
// @Router	/user/{id} [patch]
func (u *User) Patch(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, errors.New("id is not valid"))
		return
	}
	var user request.User
	err = ctx.ShouldBindJSON(&user)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	err = service.User.Patch(ctx, id, &model.User{
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Phone:    user.Phone,
		Email:    user.Email,
		Locked:   false,
	})
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Accepted(ctx)
}

// Delete godoc
//
// @Summary	Delete user by id
// @Tags		User
// @Param		id				path		string				true		"User id"
// @Success		200				{object}	response.Response
// @Failure		400				{object}	response.Response
// @Failure		401				{object}	response.Response
// @Failure		500				{object}	response.Response
// @Router	/user/{id} [delete]
func (u *User) Delete(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, errors.New("id is not valid"))
		return
	}
	err = service.User.Delete(ctx, id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}
