package testutil

import (
	"atm/pkg/model"
	"atm/pkg/service"
	"errors"
)

type dummyCardSvc struct {
	opts DummyCardTestOptions
}

func (d *dummyCardSvc) InsertCard(card model.Card) error {
	if d.opts.ErrOnInsert {
		return errors.New("insert card error")
	}

	return nil
}

func (d *dummyCardSvc) RemoveCard() error {
	if d.opts.ErrOnRemove {
		return errors.New("remove card error")
	}

	return nil
}

type DummyCardTestOptions struct {
	ErrOnInsert bool
	ErrOnRemove bool
}

func NewDummyCardSvc(opts DummyCardTestOptions) service.CardInterface {
	return &dummyCardSvc{opts: opts}
}
