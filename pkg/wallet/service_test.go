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
	repeatPayment,err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}
	
	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): can't find payment by id, error = %v", err)
		return
	}
	

	if !reflect.DeepEqual(savedPayment,repeatPayment){
		t.Errorf("newPayment(): error = %v", err)
		return
	}

	if savedPayment.ID != repeatPayment.ID {
		t.Errorf("Repeat: payments are equal! error = %v", err)
		return	
	}

}