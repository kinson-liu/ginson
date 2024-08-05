package core

import (
	"ginson/config"
	"github.com/rs/zerolog/log"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	gin.SetMode(config.Config.Service.Mode)
	s := initServer("0.0.0.0:"+config.Config.Service.Port, Routers())
	log.Error().Msg(s.ListenAndServe().Error())
}

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Minute
	s.WriteTimeout = 10 * time.Minute
	s.MaxHeaderBytes = 1 << 20
	return s
}
