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

// pageAnalyserHandler represents the struct for page analyser handler
type pageAnalyserHandler struct {
	analyserService service.AnalyserService
}

// NewPageAnalyserHandler creates and returns PageAnalyserHandler object
func NewPageAnalyserHandler(analyserService service.AnalyserService) PageAnalyserHandler {
	return &pageAnalyserHandler{analyserService: analyserService}
}

// Analyse handler for analyse API
func (h *pageAnalyserHandler) Analyse(c *gin.Context) {
	var pageAnalyseRequest model.PageAnalyseRequest
	if err := c.BindJSON(&pageAnalyseRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.NewError(model.ErrorCodeInvalidRequest, err.Error(), http.StatusBadRequest))
		return
	}

	response, err := h.analyserService.Analyse(pageAnalyseRequest)
	if err != nil {
		handlerError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

// handlerError handler error and custom error response
func handlerError(c *gin.Context, err error) {
	if modelErr, ok := err.(*model.Err); ok && modelErr.StatusCode != 0 {
		c.JSON(modelErr.StatusCode, err)
		c.Abort()
		return
	}
	c.JSON(http.StatusInternalServerError, err.Error())
	c.Abort()
}

