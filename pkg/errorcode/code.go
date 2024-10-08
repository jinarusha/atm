package errorcode

const (
	NoCardFound    = "no card found"
	InsertCardFail = "failed to insert card"
	RemoveCardFail = "failed to remove card"

	PinNumberCheckFail    = "failed to check card pin number"
	InvalidPinNumber      = "invalid pin number"
	PinNumberNotValidated = "pin number not validated"

	GetAccountIDsFail       = "failed to get account ids"
	NoMatchingAccountID     = "no matching account id"
	FailedToSelectAccountID = "failed to select account id"
	NoAccountSelected       = "no account selected"
	AccountIDMismatch       = "account id does not match"

	FailedToGetBalance = "failed to get balance"

	FailedToMakeDeposit = "failed to make deposit"

	IsOverdraw       = "is overdraw"
	FailedToWithdraw = "failed to withdraw"
)
