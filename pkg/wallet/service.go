package wallet

import (
	"errors"
	"fmt"

	"github.com/FrankS17/wallet/pkg/types"
	"github.com/google/uuid"
)


func New(text string) error {
	return &errorString{text}
}

type errorString struct {  			// сам тип ошибки не экспортируется
	s string
}

func (e *errorString) Error() string { 		
	return e.s
}

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments 	  []*types.Payment
	favorites	  []*types.Favorite
}

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be positive")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")


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


func (s *Service) Deposit(accountID int64, amount types.Money) error {
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
	
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}


func (s *Service) Reject(paymentID string) error {
	
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


func (s *Service) Repeat(paymentID string) (*types.Payment, error) {

	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}
	
	if account.Balance < payment.Amount {
		return nil, ErrNotEnoughBalance
	}
	account.Balance -= payment.Amount
	aaa := uuid.New().String()
	newP := &types.Payment{
		ID: aaa,
		AccountID: payment.AccountID,
		Amount: payment.Amount,
		Category: payment.Category,
		Status: types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, newP)
	return newP, nil
}


func (s *Service) FavoritePayment(paymentID, name string) (*types.Favorite, error) {
	
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	favoritePayment := &types.Favorite{
		ID: uuid.New().String(),
		AccountID: payment.AccountID,
		Name: name,
		Amount: payment.Amount,
		Category: payment.Category,
	}

	s.favorites = append(s.favorites, favoritePayment)	
	return favoritePayment, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, ErrFavoriteNotFound 
	}
	
	account, err := s.FindAccountByID(favorite.AccountID)
	if err != nil {
		return nil, ErrAccountNotFound 
	}

	fmt.Println(favorite)

	account.Balance -= favorite.Amount
	payment := &types.Payment{
		ID: uuid.New().String(),
		AccountID: favorite.AccountID,
		Amount: favorite.Amount,
		Category: favorite.Category,
		Status: types.PaymentStatusInProgress,
	}	

	s.payments = append(s.payments, payment)

	payment2 := &types.Payment{
		ID: uuid.New().String(),
		AccountID: favorite.AccountID,
		Amount: favorite.Amount,
		Category: favorite.Category,
		Status: types.PaymentStatusInProgress,
	}	
	fmt.Println(payment2)
	fmt.Println(payment2)
	fmt.Println(payment2)

	
	return payment, nil
}