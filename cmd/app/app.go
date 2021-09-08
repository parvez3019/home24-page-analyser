package app

import (
	"net/http"
	"page-analyser/cmd/config"
	"time"

	log "github.com/sirupsen/logrus"

	routerPkg "page-analyser/cmd/router"
	"page-analyser/handler"
)

type App struct {
	*config.Config
}

// NewApp takes configurations and creates an app server
func NewApp(config *config.Config) *App {
	return &App{
		Config: config,
	}
}

// Start takes router as argument and starts the app server
func (s *App) Start(router routerPkg.IRouter) {
	r := router.
		RegisterRoutes(handler.NewHandler()).
		Get()

	log.Infof("Serving on: %s", s.Port)
	server := &http.Server{
		Addr:         "0.0.0.0:" + s.Port,
		Handler:      r,
		ReadTimeout:  DefaultTimeout,
		WriteTimeout: DefaultTimeout,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Errorf("Error %s \n", err.Error())
	}

}

const (
	DefaultTimeout = 60 * time.Second
)
