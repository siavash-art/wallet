package wallet

import (
	"errors"
	"io"
	//"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"fmt"
	"github.com/google/uuid"
	"github.com/siavash-art/wallet/pkg/types"
	"sync"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")
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

	defer func() {
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

	defer func() {
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

	accounts := strings.Split(data, "|")
	accounts = accounts[:len(accounts)-1]

	for _, account := range accounts {

		value := strings.Split(account, ";")

		id, err := strconv.Atoi(value[0])
		if err != nil {
			return err
		}

		phone := types.Phone(value[1])

		balance, err := strconv.Atoi(value[2])
		if err != nil {
			return err
		}

		acc := &types.Account{
			ID:      int64(id),
			Phone:   phone,
			Balance: types.Money(balance),
		}

		s.accounts = append(s.accounts, acc)
		log.Print(account)
	}
	return nil
}

// Export all methods
func (s *Service) Export(dir string) error {
	if len(s.accounts) != 0 {
		
		accountsDir, err := os.Create(dir + "/accounts.dump")
		if err != nil {
			log.Println(err)
			return ErrFileNotFound
		}
		defer func() {
			if cerr := accountsDir.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		accList := ""

		for _, account := range s.accounts {
			ID := fmt.Sprint(account.ID) + ";"
			phone := fmt.Sprint(account.Phone) + ";"
			balance := fmt.Sprint(account.Balance)
			accList += ID
			accList += phone
			accList += balance + "\n"
		}
		_, err = accountsDir.WriteString(accList)
		if err != nil {
			return  err
		}
	}

	//export payment
	if len(s.payments) != 0 {
		
		paymentsDir, err := os.Create(dir + "/payments.dump")
		if err != nil {
			log.Println(err)
			return ErrFileNotFound
		}
		defer func() {
			if cerr := paymentsDir.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		paymentList := ""

		for _, payment := range s.payments {
			ID := fmt.Sprint(payment.ID) + ";"
			accountID := fmt.Sprint(payment.AccountID) + ";"
			amount := fmt.Sprint(payment.Amount) + ";"
			category := fmt.Sprint(payment.Category) + ";"
			status := fmt.Sprint(payment.Status)
			paymentList += ID
			paymentList += accountID
			paymentList += amount
			paymentList += category
			paymentList += status + "\n"
		}
		_, err = paymentsDir.WriteString(paymentList)
		if err != nil {
			return  err
		}
	}

	//export favorites
	if len(s.favorites) != 0 {
		favoritesDir, err := os.Create(dir + "/favorites.dump")
		if err != nil {
			log.Println(err)
			return ErrFileNotFound
		}
		defer func() {
			if cerr := favoritesDir.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()		

		favoriteList := ""

		for _, favorite := range s.favorites {
			ID := fmt.Sprint(favorite.ID) + ";"
			accountID := fmt.Sprint(favorite.AccountID) + ";"
			name := fmt.Sprint(favorite.Name) + ";"
			amount := fmt.Sprint(favorite.Amount) + ";"
			category := fmt.Sprint(favorite.Category)
			favoriteList += ID
			favoriteList += accountID
			favoriteList += name
			favoriteList += amount
			favoriteList += category + "\n"
		}
		
		_, err = favoritesDir.WriteString(favoriteList)
		if err != nil {
			return  err
		}
	}
	return nil
}

// Import all files
func (s *Service) Import(dir string) error {

	accountsFile, err := os.Open(dir + "/accounts.dump")
	if err != nil {
		log.Print(err)
		err = ErrFileNotFound
	}
	if err != ErrFileNotFound {
		defer func() {
			if cerr := accountsFile.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		content := make([]byte, 0)
		buf := make([]byte, 1024)

		for {
			read, err := accountsFile.Read(buf)
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

		accounts := strings.Split(data, "\n")
		accounts = accounts[:len(accounts)-1]

		for _, account := range accounts {
			value := strings.Split(account, ";")
			id, err := strconv.Atoi(value[0])
			if err != nil {
				return err
			}
			phone := types.Phone(value[1])
			balance, err := strconv.Atoi(value[2])
			if err != nil {
				return err
			}
			acc := &types.Account{
				ID:      int64(id),
				Phone:   phone,
				Balance: types.Money(balance),
			}

			s.accounts = append(s.accounts, acc)
			log.Print(account)
		}
	}
	//import payments.dump
	paymentFile, err := os.Open(dir + "/payments.dump")
	if err != nil {
		log.Print(err)
		err = ErrFileNotFound
	}
	if err != ErrFileNotFound {
		defer func() {
			if cerr := paymentFile.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		paymentContent := make([]byte, 0)
		paymentBuf := make([]byte, 1024)

		for {
			read, err := paymentFile.Read(paymentBuf)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Print(err)
				return ErrFileNotFound
			}
			paymentContent = append(paymentContent, paymentBuf[:read]...)
		}

		paymentData := string(paymentContent)
		payments := strings.Split(paymentData, "\n")
		payments = payments[:len(payments)-1]

		for _, payment := range payments {
			value := strings.Split(payment, ";")
			id := string(value[0])
			accountID, err := strconv.Atoi(value[1])
			if err != nil {
				return err
			}
			amount, err := strconv.Atoi(value[2])
			if err != nil {
				return err
			}
			category := string(value[3])
			status := string(value[4])

			pay := &types.Payment{
				ID:        string(id),
				AccountID: int64(accountID),
				Amount:    types.Money(amount),
				Category:  types.PaymentCategory(category),
				Status:    types.PaymentStatus(status),
			}

			s.payments = append(s.payments, pay)
			log.Print(payment)
		}
	}

	//import favorites.dump
	favoriteFile, err := os.Open(dir + "/favorites.dump")
	if err != nil {
		log.Print(err)
		err = ErrFileNotFound
	}

	if err != ErrFileNotFound {
		defer func() {
			if cerr := favoriteFile.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		favContent := make([]byte, 0)
		favBuf := make([]byte, 1024)

		for {
			read, err := favoriteFile.Read(favBuf)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Print(err)
				return ErrFileNotFound
			}
			favContent = append(favContent, favBuf[:read]...)
		}

		favData := string(favContent)
		favorites := strings.Split(favData, "\n")
		favorites = favorites[:len(favorites)-1]

		for _, favorite := range favorites {
			value := strings.Split(favorite, ";")
			id := string(value[0])
			accountID, err := strconv.Atoi(value[1])
			if err != nil {
				return err
			}
			name := string(value[2])
			amount, err := strconv.Atoi(value[3])
			if err != nil {
				return err
			}
			category := string(value[4])

			favorite := &types.Favorite{
				ID:        id,
				AccountID: int64(accountID),
				Name:      name,
				Amount:    types.Money(amount),
				Category:  types.PaymentCategory(category),
			}

			s.favorites = append(s.favorites, favorite)
			log.Print(favorite)
		}
	}
	return nil
}

// ExportAccountHistory - export account history by account Id
 func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	var payments []types.Payment
	for _, payment := range s.payments {
		if payment.AccountID == accountID {
			payments= append(payments, *payment)
		} 	
	}
	if payments == nil {		
		return nil, ErrAccountNotFound
	}
	return payments, nil
 }

 // HistoryToFiles get payments from ExportAccountHistory and add to file
 func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	
	if len(payments) > 0 {

		if len(payments) <= records {
			file, _ := os.OpenFile(dir + "/payments.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			defer func() {
				if cerr := file.Close(); cerr != nil {
					log.Print(cerr)
				}
			}()
						
			paymentList := ""

			for _, payment := range payments {
				ID := fmt.Sprint(payment.ID) + ";"
				accountID := fmt.Sprint(payment.AccountID) + ";"
				amount := fmt.Sprint(payment.Amount) + ";"
				category := fmt.Sprint(payment.Category) + ";"
				status := fmt.Sprint(payment.Status)
				paymentList += ID
				paymentList += accountID
				paymentList += amount
				paymentList += category
				paymentList += status + "\n"
			}
			_, err := file.WriteString(paymentList)
			if err != nil {
				return  err
			}
		
		} else {

			paymentList := ""
			counter := 0
			nextFile := 1
			var file *os.File

			for _, payment := range payments {
				if counter == 0 {
					file, _ = os.OpenFile(dir + "/payments"+fmt.Sprint(nextFile)+".dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
				}
				counter++

				paymentList = fmt.Sprint(payment.ID) + ";" + fmt.Sprint(payment.AccountID) + ";" + fmt.Sprint(payment.Amount) + ";" + fmt.Sprint(payment.Category) + ";" + fmt.Sprint(payment.Status) + "\n"
			
				_, err := file.WriteString(paymentList)
				if err != nil {
					return  err
				}
				if counter == records {
					paymentList = ""
					nextFile++
					counter =0
					file.Close()
				}
			}
		}

	}
	return nil
 }
 
 //SumPayments  return sum of payments	
 func (s *Service) SumPayments(goroutines int) types.Money {	
	wg := sync.WaitGroup{}	
	mu := sync.Mutex{}	
	sum := int64(0)
	count := 0
	i := 0
	
	if goroutines == 0 {
		count = len(s.payments) 
	} else {
		count = int(len(s.payments) / goroutines)
	}
	for i = 0; i < goroutines-1; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()			
			val := int64(0)
			payments := s.payments[index*count : (index+1)*count]
			
			for _, payment := range payments {
				val += int64(payment.Amount)
			}

			mu.Lock()
			sum += val
			mu.Unlock()

		}(i)
	} 
	wg.Add(1)
	go func (){
		defer wg.Done()
		val := int64(0)
		payments := s.payments[i*count:]
		for _, payment := range payments {
			val += int64(payment.Amount)
		}
		mu.Lock()
		sum += val
		mu.Unlock()
	}()
	
	wg.Wait()

	return types.Money(sum)
 } 

 // FilterPayments filtered payments
 func (s *Service) FilterPayments(accountID int64, goroutines int) (filtPayments []types.Payment, err error) {
	
	if goroutines < 2 {
		for _, payment := range s.payments {
			if payment.AccountID == accountID {
				filtPayments = append(filtPayments, *payment)
			}
			if filtPayments == nil {
				return nil, ErrAccountNotFound
			} 
			return
		}	
	}
	wg := sync.WaitGroup{}	
	mu := sync.Mutex{}	
	max := 0
	count := len(s.payments) / goroutines
	
	for i := 0; i < goroutines; i++ {
		max += count
		wg.Add(1)
		go func(val int){
			defer wg.Done()
			sum := []types.Payment{}
			for _, payment := range s.payments {
				if payment.AccountID == accountID {
					sum = append(sum, *payment)
				}
			}
			mu.Lock()
			filtPayments = append(filtPayments, sum...)
			mu.Unlock()
		}(max)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		sum := []types.Payment{}
		for _, payment := range s.payments[max:] {
			if payment.AccountID == accountID {
				sum = append(sum, *payment)
			}
		}
		mu.Lock()
		filtPayments = append(filtPayments, sum...)
		mu.Unlock()
	}()
	wg.Wait()
	if filtPayments == nil {
		return nil, ErrAccountNotFound
	} 	
	return 
 } 
