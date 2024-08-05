package main

import (
	"ginson/core"
)

// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description JWT Token
// @type apiKey
func main() {
	core.InitConfig()
	core.InitLogger()
	core.InitDB()
	core.InitRedis()
	core.RunServer()
}
