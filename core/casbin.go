package core

import (
	"ginson/config"
	"ginson/core/const/casbin"
	"ginson/core/const/db"
	casbinEnforcer "github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/rs/zerolog/log"
)

func InitCasbin() {
	ac := &mongodbadapter.AdapterConfig{
		DatabaseName:   config.Config.DB.Schema,
		CollectionName: "casbin",
	}
	adapter, err := mongodbadapter.NewAdapterByDB(db.Mongo, ac)
	if err != nil {
		log.Fatal().Msgf("connect casbin db failed: %v", err)
	}
	model, err := casbinModel.NewModelFromString(casbin.Model)
	if err != nil {
		log.Fatal().Msgf("parse casbin model failed: %v", err)
		return
	}
	casbin.Enforcer, err = casbinEnforcer.NewEnforcer(model, adapter)
	if err != nil {
		log.Fatal().Msgf("init casbin enforcer failed: %v", err)
		return
	}
	err = casbin.Enforcer.LoadPolicy()
	if err != nil {
		log.Fatal().Msgf("load casbin policy failed: %v", err)
		return
	}
}
