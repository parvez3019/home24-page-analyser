package handler

import (
	"github.com/gin-gonic/gin"
	"home24-page-analyser/model"
	"home24-page-analyser/service"
	"net/http"
)

// PageAnalyserHandler handler interface for pageAnalyser
type PageAnalyserHandler interface {
	Analyse(c *gin.Context)
}

// pageAnalyserHandler represents the struct for page anaylser handelr
type pageAnalyserHandler struct {
	analyserService service.AnalyserService
}

// NewPageAnalyserHandler creates and returns PageAnalyserHandler object
func NewPageAnalyserHandler(analyserService service.AnalyserService) PageAnalyserHandler {
	return &pageAnalyserHandler{
		analyserService: analyserService,
	}
}

// Analyse handler for analyse API
func (h *pageAnalyserHandler) Analyse(c *gin.Context) {
	var pageAnalyseRequest model.PageAnalyseRequest
	if err := c.BindJSON(&pageAnalyseRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.analyserService.Analyse(pageAnalyseRequest)
	if err != nil {
		handlerError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func handlerError(c *gin.Context, err error) {
	c.IndentedJSON(http.StatusInternalServerError, err.Error())
}
