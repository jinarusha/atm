package context

type AtmContextKey string

const (
	AccountID         = AtmContextKey("account_id")
	CardInserted      = AtmContextKey("card_inserted")
	CardHolderName    = AtmContextKey("card_holder_name")
	CardNumber        = AtmContextKey("card_number")
	PinNumIsValidated = AtmContextKey("pin_num_is_validated")
)
