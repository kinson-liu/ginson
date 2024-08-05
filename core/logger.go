package core

import (
	"ginson/config"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func InitLogger() {
	level, err := zerolog.ParseLevel(config.Config.Log.Level)
	if err != nil {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
	rotateLogger := &lumberjack.Logger{
		Filename:   config.Config.Log.Path,
		MaxSize:    config.Config.Log.MaxSize,
		MaxAge:     config.Config.Log.MaxAge,
		MaxBackups: config.Config.Log.MaxBackups,
		LocalTime:  true,
		Compress:   true,
	}
	consoleLogger := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = log.Output(zerolog.MultiLevelWriter(consoleLogger, rotateLogger))
}
