package wallet

import (
	"fmt"
	"reflect"
	"testing"

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
func TestService_FindPaymentID_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentID(): error = %v", err)
		return
	}

	//сравниваем платежи
	if !reflect.DeepEqual(payment, got)  {
		t.Errorf("FindPaymentID(): wrong payment was returned = %v", err)
		return
	}
}

func TestService_FindPaymentID_fail(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentID(): error = %v", err)
		return
	}

	//сравниваем платежи
	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentID(): must return ErrPaymentNotFound, returned= %v", err)
		return
	}	
}

func TestService_Reject_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't change, error = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(): balance didn't change, error = %v", savedAccount)
		return
	}
} 


func TestService_Repeat_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
		
	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): can't find payment by id, error = %v", err)
		return
	}

	repeatPayment,err := s.Repeat(savedPayment.ID)
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}
	

	_, err = s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Repeat(): can't find account by id, error = %v", err)
		return
	}


	if savedPayment.ID == repeatPayment.ID {
		t.Errorf("Repeat: payments are equal! error = %v", err)
		return	
	}

}

func TestService_FindFavoriteID_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
	
	addFavoritePayment, err := s.FavoritePayment(payment.ID, "a")
	if err != nil {
		t.Errorf("FindFavoriteID(): error = %v", err)
		return
	} 

	findFavorite, err := s.FindFavoriteByID(addFavoritePayment.ID)
	if err != nil {
		t.Errorf("FindFavoriteID(): favorite payment was not found, error = %v", err)
		return
	} 

	//сравниваем платежи
	if !reflect.DeepEqual(addFavoritePayment, findFavorite)  {
		t.Errorf("FindPaymentID(): favorite payments are not equal = %v", err)
		return
	}
}

func TestService_FindFavoriteID_fail(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	_, err = s.FindFavoriteByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindFavoriteByID(): error = %v", err)
		return
	}

	//сравниваем платежи
	if err != ErrFavoriteNotFound {
		t.Errorf("FindFavoriteByID(): must return ErrPaymentNotFound, returned= %v", err)
		return
	}	
}

func TestService_FavoritePayment_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
	
	findPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Favorite(): can't find payment by id, error = %v", err)
		return
	}

	// добавляем платеж в favorite
	addFavPay,err := s.FavoritePayment(findPayment.ID,"aaa")
	if err != nil {
		t.Errorf("Favorite(): error = %v", err)
		return
	}	
	
	// находим этот добавленный платеж в слайсе favorites 
	findFavoritePayment, err := s.FindFavoriteByID(findPayment.ID)
	if err != nil {
		t.Errorf("Favorite(): error = %v", err)
		return
	}	

	// сравниваем полученный платеж из слайса favorites - с добавленным 
	if !reflect.DeepEqual(addFavPay,findFavoritePayment) {
		t.Errorf("Favorite(): error = %v", err)
		return
	}
	
}


func TestService_PayFromFavorite_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
	
	findPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): can't find payment by id, error = %v", err)
		return
	}

	// добавляем платеж в favorite
	addFavPay,err := s.FavoritePayment(findPayment.ID,"aaa")
	if err != nil {
		t.Errorf("PayFromFavorite(): error = %v", err)
		return
	}	
	
	// находим этот добавленный платеж в слайсе favorites 
	findFavoritePayment, err := s.FindFavoriteByID(findPayment.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): error = %v", err)
		return
	}	

	// сравниваем полученный платеж из слайса favorites - с добавленным 
	if !reflect.DeepEqual(addFavPay,findFavoritePayment) {
		t.Errorf("PayFromFavorite(): error = %v", err)
		return
	}

	
	newPayment, err := s.PayFromFavorite(addFavPay.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): error = %v", err)
		return
	}

	_, err = s.FindAccountByID(newPayment.AccountID)
	if err != nil {
		t.Errorf("PayFromFavorite(): account was not found, error = %v", err)
		return
	}


	
}