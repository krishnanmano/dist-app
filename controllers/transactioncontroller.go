package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"

	"dist-app/model"
	"dist-app/service"
)

type TransactionController struct {
	service service.ITransactionService
}

func NewTransactionController(srvs service.ITransactionService) *TransactionController {
	return &TransactionController{
		service: srvs,
	}
}

func (txnc TransactionController) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

func (txnc TransactionController) GetTransactions(c *gin.Context) {
	messages, err := txnc.service.FindAllTransactions(c)
	if err != nil {
		model.NewErrorMsg(http.StatusInternalServerError, "key not found", err.Error())
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (txnc TransactionController) FindTransactionByID(c *gin.Context) {
	txnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		model.NewErrorMsg(http.StatusBadRequest, "invalid task id", err.Error())
	}

	txn, err := txnc.service.FindTransactionByID(c, txnID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "key not found",
		})
		return
	}

	c.JSON(http.StatusOK, txn)
}

func (txnc TransactionController) UpdateTransactionByID(c *gin.Context) {
	txnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		model.NewErrorMsg(http.StatusBadRequest, "invalid task id", err.Error())
	}

	var incomingTxn model.Transaction
	if err := c.ShouldBind(&incomingTxn); err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewErrorMsg(http.StatusBadRequest, "Invalid input", err.Error()),
		)
		return
	}

	_, err = txnc.service.FindTransactionByID(c, txnID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "key not found",
		})
		return
	}

	savedTxn, err := txnc.service.UpdateTransactionByID(c, txnID, &incomingTxn)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			model.NewErrorMsg(http.StatusBadRequest, "internal server error", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, savedTxn)
}

func (txnc TransactionController) DeleteTransactionByID(c *gin.Context) {
	txnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		model.NewErrorMsg(http.StatusBadRequest, "invalid task id", err.Error())
	}

	err = txnc.service.DeleteTransactionByID(c, txnID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "key not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (txnc TransactionController) SaveTransaction(c *gin.Context) {
	var incomingTxn model.Transaction
	if err := c.ShouldBind(&incomingTxn); err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.NewErrorMsg(http.StatusBadRequest, "Invalid input", err.Error()),
		)
		return
	}

	savedTxn, err := txnc.service.SaveTransaction(c, &incomingTxn)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			model.NewErrorMsg(http.StatusBadRequest, "internal server error", err.Error()),
		)
		return
	}

	c.JSON(http.StatusCreated, savedTxn)
}
