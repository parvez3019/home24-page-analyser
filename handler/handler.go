package handler

import (
	"github.com/gin-gonic/gin"
	"home24-page-analyser/service"
)

type PageAnalyserHandler interface {
	Analyse(c *gin.Context)
}

type pageAnalyserHandler struct {
	analyserService service.AnalyserService
}

func NewPageAnalyserHandler(analyserService service.AnalyserService) PageAnalyserHandler {
	return &pageAnalyserHandler{
		analyserService: analyserService,
	}
}
