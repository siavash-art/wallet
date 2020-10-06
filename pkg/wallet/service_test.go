package wallet

import (
	"github.com/siavash-art/wallet/pkg/types"
	"testing"
	"fmt"
	"reflect"
	"github.com/google/uuid"
)

type testService struct {
	*Service
}
type testAccount struct {
	phone types.Phone
	balance types.Money
	payments []struct {
		amount types.Money
		category types.PaymentCategory
	}
}
var defaultTestAccount = testAccount{
	phone: "+992938638676",
	balance: 10_000_00,
	payments: []struct {
		amount types.Money
		category types.PaymentCategory
	}{
	{amount: 1_000_00, category: "cat"},
},
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit account, error = %v", err)
	}
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, errror = %v", err)
		}
	}
	return account, payments, nil
}

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

func TestService_Repeat_success(t *testing.T) {
	s := Service{}
	s.RegisterAccount("+992938638676")

	account, err := s.FindAccountByID(1)
	if err != nil {
		t.Errorf("\n FindAccountByID(): error = %v", err)
	}

	err = s.Deposit(account.ID, 200_00)
	if err != nil {
		t.Errorf("\n Deposit(): error = %v", err)
	}

	payment, err := s.Pay(account.ID, 100_00, "cat")
	if err != nil {
		t.Errorf("\n Deposit(): error = %v", err)
	}

	pay, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("\n Deposit(): error = %v", err)
	}

	pay, err = s.Repeat(pay.ID)
	if err != nil {
		t.Errorf("Repeat(): Error(): can't pay (%v): %v", pay.ID, err)
	}
}


func TestService_FavoritePayment_success(t *testing.T) {
	
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payment := payments[0]
	got, err := s.FavoritePayment(payment.ID, "auto")
	if err != nil {
		t.Errorf("FavoriteFromPayment(): error = %v", err)
		return
	}
	
	
	if reflect.DeepEqual(payment, got) {
		t.Errorf("FavoritePayment(): wrong payment returned = %v", err)
		return
	}
}
func TestService_FavoritePayment_fail(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payment := payments[0]
	got, err := s.FavoritePayment(payment.ID, "auto")
	if err != nil {
		t.Errorf("FavoritePayment(): error = %v", err)
		return
	}
	
	
	if reflect.DeepEqual(payment, got) {
		t.Errorf("FavoritePayment(): wrong payment returned = %v", err)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	
	_, err = s.PayFromFavorite(uuid.New().String())
	if err == nil {
		t.Error("PayFromFavorite(): must return error, returned nil")
		return
	}
	
	
	if err == ErrPaymentNotFound {
		t.Errorf("PayFromFavorite(): must return ErrFavoriteNotFound, returned = %v", err)
		return
	}
}

func TestService_FindFavoriteByID_success(t *testing.T) {
	
	s := newTestService()
	_, favorites, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	favorite := favorites[0]
	got, err := s.FindFavoriteByID(favorite.ID)
	if err == nil {
		t.Errorf("FindFavoriteByID(): error = %v", err)
		return
	}
	
	
	if reflect.DeepEqual(favorite, got) {
		t.Errorf("FindFavoriteByID(): wrong payment returned = %v", err)
		return
	}
}
func TestService_FindFavoriteByID_fail(t *testing.T) {
	
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	
	_, err = s.FindFavoriteByID(uuid.New().String())
	if err == nil {
		t.Error("FindFavoriteByID(): must return error, returned nil")
		return
	}
	
	
	if err != ErrFavoriteNotFound {
		t.Errorf("FindFavoriteByID(): must return ErrFavoriteNotFound, returned = %v", err)
		return
	}
}

func TestService_Favorite_success_user(t *testing.T) {
	svc := Service{}
	account, err := svc.RegisterAccount("+9920000000001")
	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}
	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}
	payment, err := svc.Pay(account.ID, 10_00, "auto")
	if err != nil {
		t.Errorf("Pay() Error() can't pay for an account(%v):%v", account, err)
		}
		
	favorite, err := svc.FavoritePayment(payment.ID, "megafon")
		if err != nil {
			t.Errorf("FavoritePayment() Error() can't for and favorite(%v): %v", favorite, err)
	}	
	paymentFavorite, err := svc.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("PayFromFavorite() Error() can't for an favorite(%v): %v", paymentFavorite, err)
	}
}
