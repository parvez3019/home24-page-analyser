package router

import (
	"net/http"
	"net/http/httputil"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"home24-page-analyser/cmd/config"
	"home24-page-analyser/handler"
)

// IRouter represents an interface to router
type IRouter interface {
	RegisterRoutes(handler handler.PageAnalyserHandler) *Router
	Get() *gin.Engine
}

// Router represents abstraction over gin router with server configs
type Router struct {
	*config.Config
	*gin.Engine
}

// NewRouter takes configs as params and returns router
func NewRouter(config *config.Config) *Router {
	log.Info("Setting up endpoints...")
	router := gin.New()

	router.Use(recovery())

	return &Router{Config: config, Engine: router}
}

// RegisterRoutes take server context and client as params and register application routes under routing groups
func (r *Router) RegisterRoutes(handler handler.PageAnalyserHandler) *Router {
	r.POST("/analyse", handler.Analyse)
	return r
}

// Get return gin engine instance
func (r *Router) Get() *gin.Engine {
	return r.Engine
}

// recovery panic recovery and logging
func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, true)
				log.Printf("[Recovery] panic recovered: \n %s \n %s \n", string(httpRequest), err)
				log.Printf("[Recovery] Stack Trace: %s\n", string(debug.Stack()))
				c.JSON(http.StatusInternalServerError, err)
			}
		}()
		c.Next()
	}
}
