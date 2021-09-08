package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	Analyse(c *gin.Context)
}

type handler struct{}

func NewHandler() Handler {
	return &handler{}
}
