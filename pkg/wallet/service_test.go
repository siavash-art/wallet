package wallet

import (
	"log"
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

func TestService_ExportToFile_success(t *testing.T) {
	
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}	
	err = s.ExportToFile("../data/export.txt")
	if err == nil {
		t.Error("returned nil")
		return
	}	
}

func TestService_ImportFromFile_success(t *testing.T) {
	var svc Service
	svc.RegisterAccount("+992938638676")
	svc.RegisterAccount("+992938638677")
	svc.Deposit(1, 100_00)
	svc.Pay(1, 50_00, "cat")
	svc.Deposit(2, 100_00)
	svc.Pay(2, 50_00, "food")
	err := svc.Import("../data")
	if err != nil {
		t.Error("retur error func Import")
	return
	}	
}

func TestService_Export_success(t *testing.T) {
	var svc Service
	svc.RegisterAccount("+992938638676")
	svc.RegisterAccount("+992938638677")
	svc.Deposit(1, 100_00)
	svc.Pay(1, 50_00, "cat")
	svc.Deposit(2, 100_00)
	svc.Pay(2, 50_00, "food")
	err := svc.Export("../../data")
	if err != nil {
		t.Error("retur error func Export")
	return
	}	
}

func TestService_Export_fail(t *testing.T) {
	var svc Service
	svc.RegisterAccount("+992938638676")
	svc.RegisterAccount("+992938638677")
	err := svc.Export("../../data")
	if err != nil {
		t.Error("retur error func Export")
	return
	}	
}

func TestService_import_success(t *testing.T) {
	var svc Service
	svc.RegisterAccount("+992938638676")
	svc.RegisterAccount("+992938638677")
	svc.Deposit(1, 100_00)
	svc.Pay(1, 50_00, "cat")
	svc.Deposit(2, 100_00)
	svc.Pay(2, 50_00, "food")
	err := svc.Import("../../data")
	if err != nil {
		t.Error("retur error func Export")
	return
	}	
}

func TestService_ExportAccountHistory_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+992938638676")
	svc.Deposit(1, 100_00)
	
	_, err := svc.Pay(1, 10_00, "auto")
	if err != nil {
		fmt.Println(err)
		return
	}
	account, err := svc.ExportAccountHistory(1)
	if err != nil {
		t.Errorf("\n expected: %v \n result: %v", account, err)
	}
}

func TestService_ExportAccountHistory_notFound(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+992938638676")

	account, err := svc.ExportAccountHistory(11)
	if err == nil {
		t.Errorf("\n expected: %v \n result: %v", account, err)
	}
}

func TestService_HistoryToFiles_notFound(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+992938638676")

	account, err := svc.ExportAccountHistory(11)
	if err == nil {
		t.Errorf("\n expected: %v \n result: %v", account, err)
	}
}

func TestService_HistoryToFiles_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+992938638676")
	svc.Deposit(1, 100_00)
	
	_, err := svc.Pay(1, 10_00, "auto")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	payments, err := svc.ExportAccountHistory(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	svc.HistoryToFiles(payments, "../../data", 2)
}

func TestService_SumPayments_success(t *testing.T) {
	var s Service

	account, err := s.RegisterAccount("+992938638676")
	if err != nil {
		t.Errorf("error registration account, account = %v", account)
	}

	err = s.Deposit(account.ID, 300)
	if err != nil {
		t.Errorf("error Depposit, error = %v", err)
	}

	_, err = s.Pay(account.ID, 10, "food")
	_, err = s.Pay(account.ID, 20, "sport")
	_, err = s.Pay(account.ID, 30, "food")
	_, err = s.Pay(account.ID, 40, "food")
	_, err = s.Pay(account.ID, 50, "sport")
	_, err = s.Pay(account.ID, 60, "food")

	if err != nil {
		t.Errorf("error Pay, error = %v", err)
	}

	want := types.Money(210)
	got := s.SumPayments(2)

	if want != got {
		t.Errorf("error SumPayments, want = %v, got = %v", want, got)
	}
}

func BenchmarkSumPayments(b *testing.B) {
	var s Service

	account, err := s.RegisterAccount("+992938638676")
	if err != nil {
		b.Errorf("error registration account, account = %v", account)
	}

	err = s.Deposit(account.ID, 300)
	if err != nil {
		b.Errorf("error Depposit, error = %v", err)
	}

	_, err = s.Pay(account.ID, 10, "food")
	_, err = s.Pay(account.ID, 20, "sport")
	_, err = s.Pay(account.ID, 30, "food")
	_, err = s.Pay(account.ID, 40, "food")
	_, err = s.Pay(account.ID, 50, "sport")
	_, err = s.Pay(account.ID, 60, "food")

	if err != nil {
		b.Errorf("error Pay, error = %v", err)
	}

	want := types.Money(210)
	got := s.SumPayments(2)

	if want != got {
		b.Errorf("error SumPayments, want = %v, got = %v", want, got)
	}
} 

func BenchmarkFilterPayments(b *testing.B) {
	svc := &Service{}
	account, err := svc.RegisterAccount("+992938638676")
	if err != nil {
		b.Errorf("error account = %v", err)
	}
	svc.Deposit(account.ID, 200)
	svc.Pay(account.ID, 10, "auto")
	svc.Pay(account.ID, 20, "food")
	svc.Pay(account.ID, 30, "food")
	svc.Pay(account.ID, 40, "food")
	svc.Pay(account.ID, 50, "food")
	svc.Pay(account.ID, 60, "food")

	filt, err := svc.FilterPayments(account.ID, 2)
	if err != nil {
		b.Errorf("error FilterPayments = %v", err)
	}
	log.Println(filt)
}
