package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending string = "pending"
	TransactionCompleted string = "completed"
	TransactionError string = "error"
	TransctionConfirmed string = "confirmed"
)

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transaction []Transaction
}

type Transaction struct {
	Base `valid:"required"`
	AccountFrom *Account `valid:"-"`
	AccountFromID string `gorm:"column:account_from_id;type:uuid" valid:"notnull"`
	Amount float64 `json:"amount" gorm:"type:float" valid:"notnull"`
	PixKeyTo *PixKey `valid:"-"`
	PixKeyIdTo string `gorm:"column:pix_key_id_to;type:uuid" valid:"notnull"`
	Status string `json:"status" gorm:"type:varchar(20)" valid:"notnull"`
	Description string `json:"description" gorm:"type:varchar(255)" valid:"-"`
	CancelDescription string `json:"cancel_description" gorm:"type:varchar(255)" valid:"-"`
}

func (t *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(t)
	
	if t.Amount <= 0 {
		return errors.New("The amount must be greater than 0")
	}
	
	if t.Status!= TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionError {
		return errors.New("invalid status")
	}

	if t.PixKeyTo.AccountID == t.AccountFrom.ID {
		return errors.New("the souce and destination account cannot be the same")
	}

	if err!= nil {
    return err
  }
	return nil
}

func NewTransaction(AccountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: AccountFrom,
    Amount: amount,
    PixKeyTo: pixKeyTo,
    Status: TransactionPending,
    Description: description,
	}

	transaction.ID = string(uuid.NewV4().String())
	transaction.CreatedAt = time.Now()

	err:= transaction.isValid()

	if err!= nil {
    return nil, err
  }

	return &transaction, nil
}

func (t *Transaction) Complete() error {
	t.Status = TransactionCompleted
  t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}
func (t *Transaction) Confirm() error {
	t.Status = TransctionConfirmed
  t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Cancel(description string) error {
	t.Status = TransactionError
  t.UpdatedAt = time.Now()
	t.Description = description
	err := t.isValid()
	return err
}