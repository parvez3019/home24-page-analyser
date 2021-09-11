package app

import (
	"home24-page-analyser/cmd/config"
	httpWrapper "home24-page-analyser/http"
	"home24-page-analyser/service"
	"home24-page-analyser/service/html_parser"
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

// StartServer takes router as argument and starts the app server
func (s *App) StartServer(router routerPkg.IRouter) {
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

// resolveDependencies resolve dependencies via dependency injection
func resolveDependencies() handler.PageAnalyserHandler {
	httpWrapperClient := httpWrapper.NewClientWrapper(
		&http.Client{Timeout: DefaultTimeout},
	)
	analyserService := service.NewAnalyzerService(httpWrapperClient, html_parser.NewParser())
	return handler.NewPageAnalyserHandler(analyserService)
}

const (
	DefaultTimeout = 60 * time.Second
)
