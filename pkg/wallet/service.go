package wallet

import (	
	"io"
	"github.com/siavash-art/wallet/pkg/types"
	"github.com/google/uuid"
	"errors"
	"log"
	"os"
	"strconv"	
	"strings"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("payment not found")
var ErrFileNotFound = errors.New("file not found")

// Service payments of accounts
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

// RegisterAccount asdasd asdasd
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

// Deposit balance
func (s *Service) Deposit(AccountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == AccountID {
			account = acc
			break
		}
	}
	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

// Pay users payments
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount

	paymentID := uuid.New().String()

	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

// FindAccountByID find account by id
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {

	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

// FindPaymentByID find payment by account id
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {

	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

// Reject changes the payment status to PaymentStatusFail
func (s *Service) Reject(paymentID string) error {

	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return err
	}

	account, err := s.FindAccountByID(payment.AccountID)

	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount

	return nil
}

// Repeat repeat payment
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {

	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	pay, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}

	return pay, nil
}

//FavoritePayment adddddd
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {

	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return nil, err
	}

	genID := uuid.New().String()

	newFavorite := &types.Favorite{
		ID:        genID,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, newFavorite)

	return newFavorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {

	for _, favorite := range s.favorites {

		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound

}

// PayFromFavorite pay from favorite
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {

	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// ExportToFile exports accounts to file
func (s *Service) ExportToFile(path string) error {

	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}

	defer func(){
		if err2 := file.Close(); err2 != nil {
			log.Print(err2)
		}
	}()

	list := ""

	for _, account := range s.accounts {
		ID := strconv.Itoa(int(account.ID)) + ";"
		phone := string(account.Phone) + ";"
		balance := strconv.Itoa(int(account.Balance))

		list += ID
		list += phone
		list += balance + "|"
	}

	_, err = file.Write([]byte(list))

	if err != nil {
		log.Print()
		return ErrFileNotFound
	}

	return nil
}

// ImportFromFile import accounts from file
func (s *Service) ImportFromFile(path string) error {
	
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	
	defer func(){
		if err2 := file.Close(); err2 != nil {
			log.Print(err2)
		}
	}()	
	
	content := make([]byte, 0)	
	buf := make([]byte, 4)
	
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Print(err)
			return ErrFileNotFound
		}
		content = append(content, buf[:read]...)
	}

	data := string(content)
	
	accounts := strings.Split(data, ":")
	accounts = accounts[:len(accounts)-1]
	
	for _, account := range accounts {
		
		value := strings.Split(account, "/")
		
		id,err := strconv.Atoi(value[0])
		if err!=nil {
			return err
		}

		phone :=types.Phone(value[2])
		
		balance, err := strconv.Atoi(value[1])
		if err!=nil {
			return err
		}
		
		acc := &types.Account {
			ID: int64(id),
			Phone: phone,
			Balance: types.Money(balance),
		}

		s.accounts = append(s.accounts, acc)
		log.Print(account)
	}
	return nil
}