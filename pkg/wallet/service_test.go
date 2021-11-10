package wallet

import (
	"reflect"
	"testing"

	"github.com/FrankS17/wallet/pkg/types"
	"github.com/google/uuid"
)

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
