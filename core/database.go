package core

import (
	"context"
	"ginson/config"
	"ginson/core/const/db"
	"github.com/rs/zerolog/log"
)

func InitDB() {
	switch config.Config.DB.Type {
	case "mysql":
		InitMysql()
	case "mongo":
		InitMongo()
	}
}

func reloadDB() {
	log.Info().Msg("DB reloading ......")
	if db.Mysql != nil {
		sqlDB, _ := db.Mysql.DB()
		_ = sqlDB.Close()
	}
	if db.Mongo != nil {
		_ = db.Mongo.Disconnect(context.Background())
	}
	InitDB()
	log.Info().Msg("DB reloaded.")
}
