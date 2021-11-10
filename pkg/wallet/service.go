package wallet

import (
	"errors"
	"fmt"

	"github.com/FrankS17/wallet/pkg/types"
	"github.com/google/uuid"
)


type testService struct {
	*Service		// embedding(встраивание)
}

func newTestService() *testService {
	return &testService{Service: &Service{}} // функция конструктор
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}


var defaultTestAccount = testAccount { 
	phone:		"+992900100500",
	balance: 	10_000_00,
	payments:	[]struct {
		amount		types.Money
		category	types.PaymentCategory	
	} {
		{amount: 1_000_00, category: "auto"},
	},
}

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments 	  []*types.Payment
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


func (s *testService) Deposit(accountID int64, amount types.Money) error {
	if amount <=0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	
	if account == nil {
		return ErrAccountNotFound
	}
	
	// зачисление средств пока не рассматриваем как платеж
	account.Balance += amount
	return nil
}

func (s *testService) FindAccountByID(accountID int64) (*types.Account,error) {
	//var s *Service 
	var account *types.Account
	
	for _, account = range s.accounts {
		if account.ID != accountID {
			return nil, ErrAccountNotFound
		}
	}

	return account, nil
}


func (s *testService) Pay(accountID int64, amount types.Money, category types.PaymentCategory)(*types.Payment, error) {
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


func (s *testService) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment
	
	for _, payment = range s.payments {
		if payment.ID != paymentID {
			return nil, ErrPaymentNotFound
		}
	}
	return payment, nil
}

func (s *testService) Reject(paymentID string) error {
	
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return ErrPaymentNotFound
	}
	
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return ErrAccountNotFound
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount

	return nil
}



func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	//региструем пользователя
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}
	
	// пополняем счет
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	// выполняем платежи
	// можем создать слайс сразу нужной длины, поскольку знаем размер
	payments := make([]*types.Payment, len(data.payments))
	
	for i, payment := range data.payments {
	// тогда здесь работаем через индекс, а не через append
	payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
	}
	
	return account, payments, nil
}


func (s *testService) Repeat(paymentID string) (*types.Payment, error) {

	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	if payment.Amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		account = acc
		break
	}


	account.Balance -= payment.Amount
	payment.ID = uuid.New().String()
	s.payments = append(s.payments, payment)

	return payment, nil
}