package wallet

import (
	"fmt"
	"testing"
)

func TestService_FindAccountByID_success(t *testing.T) {
	
	svc := &Service{}
	account, err := svc.RegisterAccount("+992900708090")
	if err != nil {
	fmt.Println(err)
	return
	}	

	accountID, err := svc.FindAccountByID(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	
	if account.ID != accountID.ID {
		t.Errorf("invalid result, expected: %v, actual: %v",account.ID,accountID.ID)
	}
}

func TestService_FindAccountById_notFound(t *testing.T) {
	

	svc := &Service{}
	account, err := svc.RegisterAccount("+992900708090")
	if err != nil {
	fmt.Println(err)
	return
	}

	

	accountID, err := svc.FindAccountByID(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	if accountID.ID != account.ID{
		 fmt.Println(ErrAccountNotFound)
	}
	
}

func TestService_Reject_success(t *testing.T) {
	svc := &Service{}
	payment, err := svc.Pay(1,10,"auto")
	
	if err != nil {
		fmt.Println(err)
		return
	} 

	FindPayment, err := svc.FindPaymentByID(payment.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if payment != FindPayment {
		fmt.Println(ErrAccountNotFound)
	}

	err = svc.Reject(FindPayment.ID)
	if err != nil {
		fmt.Println("Wrong!")
	}

}

func TestService_Reject_notFound(t *testing.T) {
	svc := &Service{}
	payment, err := svc.Pay(1,10,"auto")
	
	if err != nil {
		fmt.Println(err)
		return
	} 

	FindPayment, err := svc.FindPaymentByID(payment.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if payment != FindPayment {
		fmt.Println(ErrAccountNotFound)
	}

	err = svc.Reject(FindPayment.ID)
	if err != nil {
		fmt.Println("Wrong!")
	}

}