package controllers

import (
	"dist-app/model"
	"dist-app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GossipController struct {
	service service.IDistAppService
}

func NewGossipController(service service.IDistAppService) *GossipController {
	return &GossipController{
		service: service,
	}
}

func (dac GossipController) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

func (dac GossipController) PublishMessage(c *gin.Context) {
	var msg model.Message
	if err := c.ShouldBind(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "invalid message passed",
		})
		return
	}

	dac.service.SaveMessage(msg)
	c.JSON(http.StatusCreated, gin.H{
		"status": "inserted",
	})
}
