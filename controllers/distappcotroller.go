package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	memlist "dist-app/memberlist"
	"dist-app/model"
	"dist-app/service"
)

type DistappController struct {
	service    service.IDistAppService
	gossipNode memlist.GossipNode
}

func NewDistappController(srvs service.IDistAppService, gossipNode memlist.GossipNode) *DistappController {
	return &DistappController{
		service:    srvs,
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
		log.Println("marshaling failed: ", err)
	}

	dac.gossipNode.HandleMessage(msgByteArr)
	c.JSON(http.StatusCreated, gin.H{
		"status": "inserted",
	})
}
