package wallet

import (
	"testing"
)

func TestService_RegisterAccount_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+992938638676")

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("\n expected: %v \n result: %v", account, err)
	}
}

func TestService_FindAccoundById_notFound(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+992938638676")

	account, err := svc.FindAccountByID(11)
	if err == nil {
		t.Errorf("\n expected: %v \n result: %v", account, err)
	}
}

func TestService_Reject_success(t *testing.T){
	
	s := Service{}
	
	s.RegisterAccount("+992938638676")

	account,err := s.FindAccountByID(1)	
	if err != nil{
		t.Errorf("\n FindAccountByID(): error = %v", err)
	}

	err = s.Deposit(account.ID, 200_00)	
	if err != nil {
		t.Errorf("\n Deposit(): error = %v", err)
	}

	payment, err := s.Pay(account.ID, 200_00, "cat")	
	if err != nil {
		t.Errorf("\n Pay(): error = %v", err)
	}

	pay, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("\n FindPaymentByID(): error = %v", err)
	}

	err = s.Reject(pay.ID)

	if err != nil {
		t.Errorf("\n Reject: error = %v", err)
	}
}
