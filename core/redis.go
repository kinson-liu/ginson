package core

import (
	"context"
	"ginson/config"
	"ginson/core/const/cache"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"strings"
)

func InitRedis() {
	switch config.Config.Cache.Type {
	case "redis":
		cache.Redis = redis.NewClient(&redis.Options{
			Addr:     config.Config.Cache.Addr,
			Password: config.Config.Cache.Password,
			DB:       config.Config.Cache.DB,
		})
	case "redis cluster":
		redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: strings.Split(config.Config.Cache.Addr, ","),
		})
	}
	_, err := cache.Redis.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal().Msgf("redis connect err: %v", err)
	}
}
func reloadRedis() {
	if cache.Redis != nil {
		_ = cache.Redis.Close()
	}
	InitRedis()
}
