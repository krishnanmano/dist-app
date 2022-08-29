package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dist-app/model"
	"dist-app/service"
)

type GossipController struct {
	txnService service.ITransactionService
}

func NewGossipController(txnSrvs service.ITransactionService) *GossipController {
	return &GossipController{
		txnService: txnSrvs,
	}
}

func (dac GossipController) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

func (dac GossipController) PublishMessage(c *gin.Context) {
	var event model.PublishEvent
	if err := c.ShouldBind(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "invalid message passed",
		})
		return
	}

	switch event.EventType {
	case model.CREATE, model.UPDATE:
		savedMsg, err := dac.txnService.SaveTransaction(c, &event.Transaction)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				model.NewErrorMsg(http.StatusBadRequest, "internal server error", err.Error()),
			)
			return
		}
		c.JSON(http.StatusOK, savedMsg)
	}

}
