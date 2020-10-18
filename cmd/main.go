package main

import (
	//"fmt"
	"github.com/siavash-art/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	// svc.RegisterAccount("+992938638676")
	// svc.RegisterAccount("+992938638677")
	// svc.Deposit(1, 100_00)
	// svc.Pay(1, 50_00, "cat")
	// svc.Deposit(2, 100_00)
	// svc.Pay(2, 50_00, "food")
	
	// payment, err := svc.Pay(1, 10_00, "auto")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }	
	// favorite, err := svc.FavoritePayment(payment.ID, "school")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// svc.PayFromFavorite(favorite.ID)
	
	// payment, err = svc.Pay(2, 10_00, "auto")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }	
	// favorite, err = svc.FavoritePayment(payment.ID, "school")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// svc.PayFromFavorite(favorite.ID)
	
	//svc.ExportToFile("../data/export.txt")
	//svc.ImportFromFile("../data/export.txt")
	//svc.Export("../data")
	svc.Import("../data")
	// account, err := svc.RegisterAccount("+992938638676")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(account)

	// err = svc.Deposit(account.ID, -10)
	// if err != nil {
	// 	switch err {
	// 	case wallet.ErrAmountMustBePositive:
	// 		fmt.Println("amount must be positive")
	// 	case wallet.ErrAccountNotFound:
	// 		fmt.Println("account not found")
	// 	}
	// 	return
	// }
	// fmt.Println(account.Balance)
}
