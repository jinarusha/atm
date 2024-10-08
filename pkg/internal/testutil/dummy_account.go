package testutil

import (
	"atm/pkg/model"
	"atm/pkg/service"
	"errors"
)

type dummyAcctSvc struct {
	opts DummyAcctTestOptions
}

func (d dummyAcctSvc) EnterPinNumber(card model.Card, number string) (bool, error) {
	if d.opts.ErrOnPinNumberEnter {
		return false, errors.New("failed to check pin number")
	}

	if d.opts.InvalidPinNumberEnter {
		return false, nil
	}

	return true, nil
}

func (d dummyAcctSvc) GetAccountIDs() ([]string, error) {
	if d.opts.ErrOnGetAccountIDs {
		return nil, errors.New("failed to get accountIDs")
	}

	return d.opts.AccountIDs, nil
}

func (d dummyAcctSvc) SelectAccountID(accountID string) error {
	if d.opts.ErrOnSelectAccountID {
		return errors.New("failed to select accountID")
	}

	return nil
}

func (d dummyAcctSvc) GetBalance(accountID string) (int, error) {
	if d.opts.ErrOnGetBalance {
		return 0, errors.New("failed to get balance")
	}

	return d.opts.GetBalanceAmt, nil
}

func (d dummyAcctSvc) MakeDeposit(accountID string, deposit int) (int, error) {
	if d.opts.ErrOnMakeDeposit {
		return 0, errors.New("failed to make deposit")
	}

	return d.opts.BalanceAfterDeposit, nil
}

func (d dummyAcctSvc) Withdraw(accountID string, withdrawAmount int) (int, error) {
	if d.opts.ErrOnWithdraw {
		return 0, errors.New("failed to withdraw")
	}

	return d.opts.BalanceAfterWithdraw, nil
}

type DummyAcctTestOptions struct {
	ErrOnPinNumberEnter   bool
	InvalidPinNumberEnter bool

	ErrOnGetAccountIDs   bool
	ErrOnSelectAccountID bool
	AccountIDs           []string

	ErrOnGetBalance bool
	GetBalanceAmt   int

	ErrOnMakeDeposit    bool
	BalanceAfterDeposit int

	ErrOnWithdraw        bool
	BalanceAfterWithdraw int
}

func NewDummyAccountSvc(opts DummyAcctTestOptions) service.AccountInterface {
	return &dummyAcctSvc{opts: opts}
}
