package wallet

import (
	"errors"

	"github.com/FrankS17/wallet/pkg/types"
	"github.com/google/uuid"
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
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")




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

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory)(*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		account = acc
		break
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID: paymentID,
		AccountID: accountID,
		Amount: amount,
		Category: category,
		Status: types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}


func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment
	
	for _, payment = range s.payments {
		if payment.ID != paymentID {
			return nil, ErrPaymentNotFound
		}
	}
	return payment, nil
}

func (s *Service) Reject(paymentID string) error {

	var amount types.Money
	
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return ErrPaymentNotFound
	}
	
	payment.Status = types.PaymentStatusFail
	
	for _, acc := range s.accounts {
		if payment.AccountID != acc.ID {
			return ErrAccountNotFound
		}
		acc.Balance += amount
	}
	return nil
}
