package service

import "atm/pkg/model"

type AccountInterface interface {
	EnterPinNumber(card model.Card, number string) (bool, error)

	GetAccountIDs() ([]string, error)
	SelectAccountID(accountID string) error

	GetBalance(accountID string) (int, error)
	MakeDeposit(accountID string, deposit int) (int, error)
	Withdraw(accountID string, withdrawAmount int) (int, error)
}
