package service

import (
	"context"
	memlist "dist-app/memberlist"
	"dist-app/model"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ITransactionService interface {
	FindAllTransactions(ctx context.Context) (*[]model.Transaction, error)
	UpdateTransactionByID(ctx context.Context, txnID uuid.UUID, txn *model.Transaction) (*model.Transaction, error)
	FindTransactionByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
	SaveTransaction(ctx context.Context, msg *model.Transaction) (*model.Transaction, error)
	DeleteTransactionByID(ctx context.Context, txnID uuid.UUID) error
}

type transactionService struct {
	model      model.ITransaction
	gossipNode *memlist.GossipNode
}

var _ ITransactionService = new(transactionService)

func NewTransactionService(gossipNode *memlist.GossipNode, dao model.ITransaction) *transactionService {
	return &transactionService{
		gossipNode: gossipNode,
		model:      dao,
	}
}

func (txnSrvs transactionService) SaveTransaction(ctx context.Context, txn *model.Transaction) (*model.Transaction, error) {
	txn.TxnID = uuid.New()
	txn.CreatedAt = time.Now().UnixMilli()
	txn.UpdatedAt = txn.CreatedAt
	return txnSrvs.model.Save(ctx, *txn)
}

func (txnSrvs transactionService) UpdateTransactionByID(ctx context.Context, txnID uuid.UUID, incomingTxn *model.Transaction) (*model.Transaction, error) {
	incomingTxn.TxnID = txnID
	existingTxn, err := txnSrvs.model.FindByID(ctx, model.Transaction{TxnID: txnID})
	if err != nil {
		return nil, err
	}

	if existingTxn.TxnID == incomingTxn.TxnID {
		existingTxn.TxnType = incomingTxn.TxnType
		existingTxn.AccountType = incomingTxn.AccountType
		existingTxn.TxnSubmitter = incomingTxn.TxnSubmitter
		existingTxn.TxnReceiver = incomingTxn.TxnReceiver
		existingTxn.TxnAmount = incomingTxn.TxnAmount
		existingTxn.Remarks = incomingTxn.Remarks
		existingTxn.UpdatedAt = time.Now().UnixMilli()
		return txnSrvs.model.Update(ctx, *existingTxn)
	}

	return nil, fmt.Errorf("transacation id not found")
}

func (txnSrvs transactionService) FindTransactionByID(ctx context.Context, txnID uuid.UUID) (*model.Transaction, error) {
	existingTxn, err := txnSrvs.model.FindByID(ctx, model.Transaction{TxnID: txnID})
	if err != nil {
		return nil, err
	}
	return existingTxn, nil
}

func (txnSrvs transactionService) FindAllTransactions(ctx context.Context) (*[]model.Transaction, error) {
	return txnSrvs.model.FindAll(ctx)
}

func (txnSrvs transactionService) DeleteTransactionByID(ctx context.Context, txnID uuid.UUID) error {
	err := txnSrvs.model.DeleteByID(ctx, model.Transaction{TxnID: txnID})
	if err != nil {
		return err
	}
	return nil
}
