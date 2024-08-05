package core

import (
	"flag"
	"fmt"
	"ginson/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

// InitConfig
// priority command -> env -> default
func InitConfig() {
	var path string
	flag.StringVar(&path, "c", os.Getenv(config.PathEnv), "choose config file.")
	flag.Parse()
	if path == "" {
		path = config.DefaultPath
	}
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&config.Config); err != nil {
			fmt.Println(err)
			return
		}
		reloadDB()
		reloadRedis()
	})
	if err = v.Unmarshal(&config.Config); err != nil {
		panic(err)
	}
	return
}
