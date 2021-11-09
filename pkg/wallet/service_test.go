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

	//expected := types.Account{1,"+992900708090",0}
	/*
	if !reflect.DeepEqual(account,accountID) {
		t.Errorf("invalid result, expected: %v, actual: %v",account,accountID)
	   }*/
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