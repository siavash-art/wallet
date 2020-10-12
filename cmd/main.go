package main

import (
	//"fmt"
	"github.com/siavash-art/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	//svc.RegisterAccount("+992938638676")
	//svc.RegisterAccount("+992938638677")
	//svc.ExportToFile("../data/export.txt")
	svc.ImportFromFile("../data/export.txt")	
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
