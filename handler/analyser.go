package handler

import (
	"github.com/gin-gonic/gin"
	"home24-page-analyser/model"
	"net/http"
)

func (h *pageAnalyserHandler) Analyse(c *gin.Context) {
	var pageAnalyseRequest model.PageAnalyseRequest
	if err := c.Bind(&pageAnalyseRequest); err != nil {
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
