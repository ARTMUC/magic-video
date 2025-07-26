package repository

import (
	"gorm.io/gorm"
)

type Tx = gorm.DB

type TransactionProvider interface {
	Transaction(fn func(tx *Tx) error) error
}

type transactionProvider struct {
	db *gorm.DB
}

func NewTransactionProvider(db *gorm.DB) TransactionProvider {
	return &transactionProvider{db: db}
}

func (t *transactionProvider) Transaction(fn func(tx *Tx) error) error {
	gormTx := t.db.Begin()
	if gormTx.Error != nil {
		return gormTx.Error
	}

	if err := fn(gormTx); err != nil {
		gormTx.Rollback()
		return err
	}

	return gormTx.Commit().Error
}
