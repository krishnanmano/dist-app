package controllers

import (
	memlist "dist-app/memberlist"
	"dist-app/model"
	"dist-app/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type DistappController struct {
	service    service.IDistAppService
	gossipNode memlist.GossipNode
}

func NewDistappController(service service.IDistAppService, gossipNode memlist.GossipNode) *DistappController {
	return &DistappController{
		service:    service,
		gossipNode: gossipNode,
	}
}

func (dac DistappController) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

func (dac DistappController) GetMessages(c *gin.Context) {
	messages := dac.service.GetMessages()
	c.JSON(http.StatusOK, messages)
}

func (dac DistappController) SaveMessage(c *gin.Context) {
	var msg model.Message
	if err := c.ShouldBind(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "invalid message passed",
		})
		return
	}

	dac.service.SaveMessage(msg)
	msgByteArr, err := json.Marshal(msg)
	if err != nil {
		log.Println("marshalling failed: ", err)
	}

	dac.gossipNode.HandleMessage(msgByteArr)
	c.JSON(http.StatusCreated, gin.H{
		"status": "inserted",
	})
}
