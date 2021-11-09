package wallet

import (
	"fmt"
	"testing"

	
)

func TestService_FindAccountById_success(t *testing.T) {
	

	svc := &Service{
	}

	
	account, err := svc.FindAccountById(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	if account.ID == 0 && err != nil  {
		t.Errorf("result wrong")
	}
	
}

func TestService_FindAccountById_notFound(t *testing.T) {
	

	svc := &Service{
	}


	account, err := svc.FindAccountById(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	if account.ID == 0 && err != nil  {
		t.Errorf("result wrong")
	}
	
}