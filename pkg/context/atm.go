package context

import "atm/pkg/model"

type AtmContext struct {
	card       *model.Card
	strValues  map[string]string
	boolValues map[string]bool
}

func NewAtmContext() *AtmContext {
	return &AtmContext{
		strValues:  make(map[string]string),
		boolValues: make(map[string]bool),
	}
}

func (ctx *AtmContext) ViewCard() *model.Card {
	return ctx.card
}

func (ctx *AtmContext) SetCard(card model.Card) {
	ctx.boolValues[string(CardInserted)] = true
	ctx.card = &card
}

func (ctx *AtmContext) HasCardInserted() bool {
	return ctx.boolValues[string(CardInserted)] && ctx.card != nil
}

func (ctx *AtmContext) SetAccountID(accountID string) {
	ctx.strValues[string(AccountID)] = accountID
}

func (ctx *AtmContext) GetAccountID() string {
	return ctx.strValues[string(AccountID)]
}

func (ctx *AtmContext) SetPinNumValid(validated bool) {
	ctx.boolValues[string(PinNumIsValidated)] = validated
}

func (ctx *AtmContext) IsPinNumValidated() bool {
	return ctx.boolValues[string(PinNumIsValidated)]
}

func (ctx *AtmContext) Clear() {
	ctx.card = nil
	ctx.strValues = make(map[string]string)
	ctx.boolValues = make(map[string]bool)
}
