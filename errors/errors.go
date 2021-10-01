package errors 

import (
	"github.com/pkg/errors"
)

var (
	ErrGeneric           	   = errors.New("something went wrong")
	ErrCreateUserFailed  	   = errors.New("failed to create user")
	ErrTransactionFailed         = errors.New("transaction failed")
	ErrInsufficientFunds 	   = errors.New("insufficient funds for the operation you're trying to perform")
)

func New(message string) error {
	return errors.New(message)
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}
