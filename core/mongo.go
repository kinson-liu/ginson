package core

import (
	"context"
	"fmt"
	"ginson/config"
	"ginson/core/const/db"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo() {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?retryWrites=true&w=majority",
		config.Config.DB.Username,
		config.Config.DB.Password,
		config.Config.DB.Host,
		config.Config.DB.Port,
	)
	opts := options.Client().ApplyURI(uri)
	log.Printf("connect mongodb success...")
	var err error
	db.Mongo, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal().Msgf("mongodb connect err: %v", err)
		return
	}
}
