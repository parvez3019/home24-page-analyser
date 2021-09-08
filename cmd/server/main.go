package main

import (
	"page-analyser/cmd/app"
	configPkg "page-analyser/cmd/config"
	"page-analyser/cmd/router"
)

func main() {
	config := configPkg.LoadConfig()
	app.NewApp(config).
		Start(router.NewRouter(config))
}
