package badgerdb

import (
	"dist-app/model"
	"encoding/json"
	"log"
)

type TransactionDAO struct{}

var _ model.ITransaction = new(TransactionDAO)

func (txnDao TransactionDAO) string(txn model.Transaction) string {
	txnInBytes, err := json.Marshal(txn)
	if err != nil {
		log.Println("Error: Failed to marshal", err)
	}
	return string(txnInBytes)
}

func (txnDao TransactionDAO) toBytes(txn model.Transaction) []byte {
	txnInBytes, err := json.Marshal(txn)
	if err != nil {
		log.Println("Error: Failed to marshal", err)
	}
	return txnInBytes
}

func (txnDao TransactionDAO) txnIDtoBytes(txn model.Transaction) []byte {
	txnIDInBytes, err := txn.TxnID.MarshalBinary()
	if err != nil {
		log.Println("Error: Failed to marshal", err)
	}
	return txnIDInBytes
}

func (txnDao *TransactionDAO) Save(txn model.Transaction) (*model.Transaction, error) {
	err := BDBClient.Set([]byte("transactions"), txnDao.txnIDtoBytes(txn), txnDao.toBytes(txn))
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func (txnDao *TransactionDAO) Update(txn model.Transaction) (*model.Transaction, error) {
	err := BDBClient.Set([]byte("transactions"), txnDao.txnIDtoBytes(txn), txnDao.toBytes(txn))
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func (txnDao *TransactionDAO) FindByID(txn model.Transaction) (*model.Transaction, error) {
	txnInBytes, err := BDBClient.Get([]byte("transactions"), txnDao.txnIDtoBytes(txn))
	if err != nil {
		return nil, err
	}
	var tempTxn model.Transaction
	json.Unmarshal(txnInBytes, &tempTxn)
	return &tempTxn, nil
}

func (txnDao *TransactionDAO) DeleteByID(txn model.Transaction) error {
	err := BDBClient.Delete([]byte("transactions"), txnDao.txnIDtoBytes(txn))
	if err != nil {
		return err
	}
	return nil
}

func (txnDao *TransactionDAO) FindAll() (*[]model.Transaction, error) {
	allTxnsInBytes, err := BDBClient.PrefixScan([]byte("transactions"))
	if err != nil {
		return nil, err
	}

	allTransaction := make([]model.Transaction, 0)
	for _, txnInBytes := range allTxnsInBytes {
		var tempTxn model.Transaction
		json.Unmarshal(txnInBytes, &tempTxn)
		allTransaction = append(allTransaction, tempTxn)
	}
	return &allTransaction, nil
}
