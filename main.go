package main

import (
	"home24-page-analyser/cmd/app"
	configPkg "home24-page-analyser/cmd/config"
	"home24-page-analyser/cmd/router"
)

func main() {
	config := configPkg.LoadConfig()
	app.NewApp(config).
		Start(router.NewRouter(config))
}
