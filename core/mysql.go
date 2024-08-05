package core

import (
	"fmt"
	"ginson/config"
	"ginson/core/const/db"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config.DB.Username,
		config.Config.DB.Password,
		config.Config.DB.Host,
		config.Config.DB.Port,
		config.Config.DB.Schema,
	)
	var err error
	if db.Mysql, err = gorm.Open(mysql.Open(dsn)); err != nil {
		log.Fatal().Msgf("mysql connect err: %v", err)
	}
	sqlDB, _ := db.Mysql.DB()
	sqlDB.SetMaxIdleConns(int(config.Config.DB.Conn.Min))
	sqlDB.SetMaxOpenConns(int(config.Config.DB.Conn.Max))
	return
}
