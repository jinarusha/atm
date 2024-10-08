package service

import (
	"atm/pkg/model"
)

type CardInterface interface {
	InsertCard(card model.Card) error
	RemoveCard() error
}
