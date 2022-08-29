package model

import (
	"context"
	"github.com/google/uuid"
)

type AccountType string

const (
	CBDS_ACCOUNT AccountType = "CBDS_ACCOUNT"
	CASH_ACCOUNT AccountType = "CASH_ACCOUNT"
)

type TransactionType string

const (
	ISSUANCE   TransactionType = "ISSUANCE"
	REDEMPTION TransactionType = "REDEMPTION"
)

type TransactionStatus string

const (
	DRAFT                TransactionStatus = "DRAFT"
	WAITING_FOR_REVIEW   TransactionStatus = "WAITING_FOR_REVIEW"
	SUBMITTED            TransactionStatus = "SUBMITTED"
	WAITING_FOR_APPROVAL TransactionStatus = "WAITING_FOR_APPROVAL"
	APPROVED             TransactionStatus = "APPROVED"
	PROCESSED            TransactionStatus = "PROCESSED"
)

const NAMESPACE_NAME = "transactions"

type Transaction struct {
	TxnID        uuid.UUID         `json:"txn_id" gorm:"primaryKey;not null;unique"`
	TxnStatus    TransactionStatus `json:"txn_status" binding:"required"`
	TxnType      TransactionType   `json:"txn_type" binding:"required"`
	AccountType  AccountType       `json:"account_type" binding:"required"`
	TxnSubmitter string            `json:"txn_submitter" binding:"required"`
	TxnReceiver  string            `json:"txn_receiver" binding:"required"`
	TxnAmount    float64           `json:"txn_amount" binding:"required"`
	Remarks      string            `json:"remarks" binding:"required"`
	CreatedAt    int64             `json:"created_at,omitempty"`
	UpdatedAt    int64             `json:"updated_at,omitempty"`
}

type ITransaction interface {
	Save(ctx context.Context, txn Transaction) (*Transaction, error)
	Update(ctx context.Context, txn Transaction) (*Transaction, error)
	FindByID(ctx context.Context, txn Transaction) (*Transaction, error)
	DeleteByID(ctx context.Context, txn Transaction) error
	FindAll(ctx context.Context) (*[]Transaction, error)
}
