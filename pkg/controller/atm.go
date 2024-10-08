package controller

import (
	"atm/pkg/context"
	"atm/pkg/errorcode"
	"atm/pkg/model"
	"atm/pkg/service"
	"errors"
	"fmt"
)

type AtmController struct {
	ctx        *context.AtmContext
	accountSvc service.AccountInterface
	cardSvc    service.CardInterface
}

type Options struct {
	accountSvc service.AccountInterface
	cardSvc    service.CardInterface
}

func NewAtmController(opts Options) *AtmController {
	return &AtmController{
		ctx:        context.NewAtmContext(),
		accountSvc: opts.accountSvc,
		cardSvc:    opts.cardSvc,
	}
}

func (ctrl *AtmController) InsertCard(card model.Card) error {
	err := ctrl.cardSvc.InsertCard(card)
	if err != nil {
		fmt.Errorf("%s: %w", errorcode.InsertCardFail, err)
		return errors.New(errorcode.InsertCardFail)
	}

	ctrl.ctx.SetCard(card)

	return nil
}

func (ctrl *AtmController) RemoveCard() error {
	if !ctrl.ctx.HasCardInserted() {
		return errors.New(errorcode.NoCardFound)
	}

	err := ctrl.cardSvc.RemoveCard()
	if err != nil {
		fmt.Errorf("%s: %w", errorcode.RemoveCardFail, err)
		return errors.New(errorcode.RemoveCardFail)
	}

	ctrl.ctx.Clear()

	return nil
}

func (ctrl *AtmController) EnterPin(pinNumber string) error {
	if !ctrl.ctx.HasCardInserted() {
		return errors.New(errorcode.NoCardFound)
	}
	card := ctrl.ctx.ViewCard()
	if card == nil {
		return errors.New(errorcode.NoCardFound)
	}

	isValid, err := ctrl.accountSvc.EnterPinNumber(*card, pinNumber)
	if err != nil {
		ctrl.ctx.SetPinNumValid(false)
		fmt.Errorf("%s: %w", errorcode.PinNumberCheckFail, err)
		return errors.New(errorcode.PinNumberCheckFail)
	}
	if !isValid {
		ctrl.ctx.Clear()
		return errors.New(errorcode.InvalidPinNumber)
	}

	ctrl.ctx.SetPinNumValid(true)

	return nil
}

func (ctrl *AtmController) GetAccountIDs() ([]string, error) {
	if !ctrl.ctx.HasCardInserted() {
		return nil, errors.New(errorcode.NoCardFound)
	}
	card := ctrl.ctx.ViewCard()
	if card == nil {
		return nil, errors.New(errorcode.NoCardFound)
	}
	if !ctrl.ctx.IsPinNumValidated() {
		return nil, errors.New(errorcode.PinNumberNotValidated)
	}

	return ctrl.accountSvc.GetAccountIDs()
}

func (ctrl *AtmController) SelectAccount(accountID string) error {
	if !ctrl.ctx.HasCardInserted() {
		return errors.New(errorcode.NoCardFound)
	}
	card := ctrl.ctx.ViewCard()
	if card == nil {
		return errors.New(errorcode.NoCardFound)
	}
	if !ctrl.ctx.IsPinNumValidated() {
		return errors.New(errorcode.PinNumberNotValidated)
	}

	accountIDs, err := ctrl.accountSvc.GetAccountIDs()
	if err != nil {
		return errors.New(errorcode.GetAccountIDsFail)
	}

	var matchFound bool
	for _, id := range accountIDs {
		if accountID == id {
			matchFound = true
			break
		}
	}

	if !matchFound {
		return errors.New(errorcode.NoMatchingAccountID)
	}

	err = ctrl.accountSvc.SelectAccountID(accountID)
	if err != nil {
		return errors.New(errorcode.FailedToSelectAccountID)
	}

	ctrl.ctx.SetAccountID(accountID)

	return nil
}

func (ctrl *AtmController) GetBalance(accountID string) (int, error) {
	if !ctrl.ctx.HasCardInserted() {
		return -1, errors.New(errorcode.NoCardFound)
	}
	card := ctrl.ctx.ViewCard()
	if card == nil {
		return -1, errors.New(errorcode.NoCardFound)
	}
	if !ctrl.ctx.IsPinNumValidated() {
		return -1, errors.New(errorcode.PinNumberNotValidated)
	}
	if ctrl.ctx.GetAccountID() == "" {
		return -1, errors.New(errorcode.NoAccountSelected)
	}
	if ctrl.ctx.GetAccountID() != accountID {
		return -1, errors.New(errorcode.AccountIDMismatch)
	}

	balance, err := ctrl.accountSvc.GetBalance(accountID)
	if err != nil {
		return -1, errors.New(errorcode.FailedToGetBalance)
	}

	return balance, nil
}

func (ctrl *AtmController) MakeDeposit(accountID string, amount int) (int, error) {
	if !ctrl.ctx.HasCardInserted() {
		return -1, errors.New(errorcode.NoCardFound)
	}
	card := ctrl.ctx.ViewCard()
	if card == nil {
		return -1, errors.New(errorcode.NoCardFound)
	}
	if !ctrl.ctx.IsPinNumValidated() {
		return -1, errors.New(errorcode.PinNumberNotValidated)
	}
	if ctrl.ctx.GetAccountID() == "" {
		return -1, errors.New(errorcode.NoAccountSelected)
	}
	if ctrl.ctx.GetAccountID() != accountID {
		return -1, errors.New(errorcode.AccountIDMismatch)
	}

	newBalance, err := ctrl.accountSvc.MakeDeposit(accountID, amount)
	if err != nil {
		return -1, errors.New(errorcode.FailedToMakeDeposit)
	}

	return newBalance, nil
}

func (ctrl *AtmController) MakeWithdrawl(accountID string, withdrawAmt int) (int, error) {
	if !ctrl.ctx.HasCardInserted() {
		return -1, errors.New(errorcode.NoCardFound)
	}
	card := ctrl.ctx.ViewCard()
	if card == nil {
		return -1, errors.New(errorcode.NoCardFound)
	}
	if !ctrl.ctx.IsPinNumValidated() {
		return -1, errors.New(errorcode.PinNumberNotValidated)
	}
	if ctrl.ctx.GetAccountID() == "" {
		return -1, errors.New(errorcode.NoAccountSelected)
	}
	if ctrl.ctx.GetAccountID() != accountID {
		return -1, errors.New(errorcode.AccountIDMismatch)
	}

	currentBalance, err := ctrl.accountSvc.GetBalance(accountID)
	if err != nil {
		return -1, errors.New(errorcode.FailedToGetBalance)
	}

	if currentBalance < withdrawAmt {
		return -1, errors.New(errorcode.IsOverdraw)
	}

	newBalance, err := ctrl.accountSvc.Withdraw(accountID, withdrawAmt)
	if err != nil {
		return -1, errors.New(errorcode.FailedToWithdraw)
	}

	return newBalance, nil
}
