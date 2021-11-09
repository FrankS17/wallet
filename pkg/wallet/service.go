package wallet

import (
	"errors"

	"github.com/FrankS17/wallet/pkg/types"
)

type Service struct {
	nextAccountId int64
	accounts      []*types.Account
	payments 	  []*types.Payment
}


func New(text string) error {
	return &errorString{text}
}

type errorString struct {  			// сам тип ошибки не экспортируется
	s string
}

func (e *errorString) Error() string { 		// но зато эекспортируется функция, которая создает ошибки этого типа
	return e.s
}

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be positive")
var ErrAccountNotFound= errors.New("account not found")



func (s *Service) FindAccountById(accountID int64) (*types.Account,error) {
	//var s *Service 

	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}