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