package wallet

import (
	"errors"

	"github.com/FrankS17/wallet/pkg/types"
)

type Service struct {
	nextAccountID int64
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



func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error){
	for _, account := range s.accounts {
	if account.Phone == phone {
		return nil, ErrPhoneRegistered
	}
}
	s.nextAccountID++
	account := &types.Account{
		ID: 		s.nextAccountID,
		Phone:		phone,
		Balance: 	0,
	}
	s.accounts = append(s.accounts,account)

	return account, nil
}



func (s *Service) FindAccountByID(accountID int64) (*types.Account,error) {
	//var s *Service 
	var account *types.Account
	
	for _, account = range s.accounts {
		if account.ID != accountID {
			return nil, ErrAccountNotFound
		}
	}

	return account, nil
}