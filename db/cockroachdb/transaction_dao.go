package cockroachdb

import (
	"context"
	"dist-app/model"
	"encoding/json"
	"fmt"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbgorm"
	"gorm.io/gorm"
	"log"
)

type transactionDAO struct{}

var _ model.ITransaction = new(transactionDAO)

func NewTransactionDAO() model.ITransaction {
	return &transactionDAO{}
}

func (txnDao transactionDAO) string(txn model.Transaction) string {
	txnInBytes, err := json.Marshal(txn)
	if err != nil {
		log.Println("Error: Failed to marshal", err)
	}
	return string(txnInBytes)
}

func (txnDao transactionDAO) toBytes(txn model.Transaction) []byte {
	txnInBytes, err := json.Marshal(txn)
	if err != nil {
		log.Println("Error: Failed to marshal", err)
	}
	return txnInBytes
}

func (txnDao transactionDAO) txnIDtoBytes(txn model.Transaction) []byte {
	txnIDInBytes, err := txn.TxnID.MarshalBinary()
	if err != nil {
		log.Println("Error: Failed to marshal", err)
	}
	return txnIDInBytes
}

func (txnDao *transactionDAO) Save(ctx context.Context, txn model.Transaction) (*model.Transaction, error) {
	if err := crdbgorm.ExecuteTx(ctx, DBClient, nil,
		func(tx *gorm.DB) error {
			result := DBClient.Create(&txn)
			if result.Error != nil {
				return result.Error
			}
			return nil
		},
	); err != nil {
		fmt.Println(err)
	}

	return &txn, nil
}

func (txnDao *transactionDAO) Update(ctx context.Context, txn model.Transaction) (*model.Transaction, error) {
	if err := crdbgorm.ExecuteTx(ctx, DBClient, nil,
		func(tx *gorm.DB) error {
			result := DBClient.WithContext(ctx).Create(&txn)
			if result.Error != nil {
				return result.Error
			}
			return nil
		},
	); err != nil {
		return nil, err
	}
	return &txn, nil
}

func (txnDao *transactionDAO) FindByID(ctx context.Context, txn model.Transaction) (*model.Transaction, error) {
	var txnFound model.Transaction
	result := DBClient.WithContext(ctx).Find(&txnFound, txn.TxnID)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("trasaction not found")
	}

	return &txnFound, nil
}

func (txnDao *transactionDAO) DeleteByID(ctx context.Context, txn model.Transaction) error {
	if err := crdbgorm.ExecuteTx(ctx, DBClient, nil,
		func(tx *gorm.DB) error {
			result := DBClient.WithContext(ctx).Delete(&model.Transaction{}, txn.TxnID)
			if result.Error != nil {
				return result.Error
			}
			return nil
		},
	); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (txnDao *transactionDAO) FindAll(ctx context.Context) (*[]model.Transaction, error) {
	var transactions []model.Transaction
	result := DBClient.Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transactions, nil
}
