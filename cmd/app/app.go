package app

import (
	"home24-page-analyser/cmd/config"
	"home24-page-analyser/service"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	routerPkg "home24-page-analyser/cmd/router"
	"home24-page-analyser/handler"
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
		RegisterRoutes(resolveDependencies()).
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

func resolveDependencies() handler.PageAnalyserHandler {
	return handler.NewPageAnalyserHandler(service.NewAnalyzerService())
}

const (
	DefaultTimeout = 60 * time.Second
)
