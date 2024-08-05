package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	apis     []Api
	validate *validator.Validate
)

type Api interface {
	Router(router *gin.RouterGroup)
	Prefix() string
}

func Init(engine *gin.Engine) {
	validate = validator.New()
	for _, api := range apis {
		api.Router(engine.Group(api.Prefix()))
	}
}
