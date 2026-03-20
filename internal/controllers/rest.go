package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-edi-document-processor/internal/logger"
	"github.com/go-edi-document-processor/internal/middleware"
)

type RestController struct {
	log *logger.Logger
	mv  *middleware.Middleware
}

func NewRestController(log *logger.Logger, mv *middleware.Middleware) *RestController {
	return &RestController{log: log, mv: mv}
}

func (c *RestController) RegisterRoutes(router *gin.Engine) {

	internal := router.Group("/internal").Use(c.mv.InternalOnly())
	{
		internal.GET("/health", c.health)
	}

	public := router.Group("/api/v1")
	{
		public.POST("/doc/send", c.sendDocument)
		public.GET("/doc/receive", c.receiveDocument)
		public.GET("/doc/:uuid", c.getDocumentByUUID)
	}
}

func (c *RestController) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (c *RestController) sendDocument(ctx *gin.Context) {
	// TODO: реализовать логику
	ctx.JSON(http.StatusOK, gin.H{
		"documents": []string{},
	})
}

func (c *RestController) receiveDocument(ctx *gin.Context) {
	// TODO: реализовать логику
	ctx.JSON(http.StatusCreated, gin.H{
		"id": "123",
	})
}

func (c *RestController) getDocumentByUUID(ctx *gin.Context) {
	id := ctx.Param("id")
	// TODO: реализовать логику
	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
