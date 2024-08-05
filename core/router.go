package core

import (
	"fmt"
	"ginson/api"
	"ginson/config"
	"ginson/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strings"
)

func Routers() (router *gin.Engine) {
	router = gin.New()
	router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		router.Use(gin.Logger())
	}

	docs.SwaggerInfo.Title = fmt.Sprintf("%s APIs", config.Config.Service.Name)
	docs.SwaggerInfo.Description = config.Config.Service.Description
	docs.SwaggerInfo.Version = config.Config.Service.Version
	if config.Config.Service.Port == "" {
		docs.SwaggerInfo.Host = config.Config.Service.Host
	} else {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.Config.Service.Host, config.Config.Service.Port)
	}
	docs.SwaggerInfo.Schemes = strings.Split(strings.TrimSpace(config.Config.Service.Scheme), ",")
	docs.SwaggerInfo.BasePath = config.Config.Service.Prefix
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	api.Init(router)
	return
}
