package controller

import (
	"atm/pkg/errorcode"
	"atm/pkg/internal/testutil"
	"atm/pkg/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInsertCard(t *testing.T) {
	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
	})
	err := ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})
	require.NoError(t, err)
	require.True(t, ctrl.ctx.HasCardInserted())

	err = ctrl.RemoveCard()
	require.NoError(t, err)
	require.False(t, ctrl.ctx.HasCardInserted())
}

func TestInsertCardError(t *testing.T) {
	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{
			ErrOnInsert: true,
		}),
	})
	err := ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})
	require.Error(t, err)
	require.False(t, ctrl.ctx.HasCardInserted())

	err = ctrl.RemoveCard()
	require.Error(t, err)
	require.False(t, ctrl.ctx.HasCardInserted())
}

func TestCardRemove(t *testing.T) {
	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
	})
	err := ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})
	require.NoError(t, err)
	require.True(t, ctrl.ctx.HasCardInserted())

	err = ctrl.RemoveCard()
	require.NoError(t, err)
	require.False(t, ctrl.ctx.HasCardInserted())
}

func TestCardRemoveError(t *testing.T) {
	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{
			ErrOnRemove: true,
		}),
	})
	err := ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})
	require.NoError(t, err)
	require.True(t, ctrl.ctx.HasCardInserted())

	err = ctrl.RemoveCard()
	require.Error(t, err)
	require.True(t, ctrl.ctx.HasCardInserted())
}

func TestPinNumber(t *testing.T) {
	ctrl := NewAtmController(Options{
		cardSvc:    testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{}),
	})
	err := ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})
	require.NoError(t, err)

	err = ctrl.EnterPin("123123231")
	require.NoError(t, err)
	require.True(t, ctrl.ctx.IsPinNumValidated())
}

func TestPinNumberSvcError(t *testing.T) {
	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			ErrOnPinNumberEnter: true,
		}),
	})
	err := ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})
	require.NoError(t, err)

	err = ctrl.EnterPin("123123231")
	require.Error(t, err)
	require.False(t, ctrl.ctx.IsPinNumValidated())
}

func TestPinNumberInvalidError(t *testing.T) {
	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			InvalidPinNumberEnter: true,
		}),
	})
	err := ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})
	require.NoError(t, err)

	err = ctrl.EnterPin("123123231")
	require.Error(t, err)
	require.EqualError(t, err, errorcode.InvalidPinNumber)
	require.False(t, ctrl.ctx.IsPinNumValidated())
}

func TestGetAccountIDs(t *testing.T) {
	expectedAccountIDs := []string{"test_account_1"}

	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			AccountIDs: expectedAccountIDs,
		}),
	})
	_ = ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})

	_ = ctrl.EnterPin("123123231")

	accountIDs, err := ctrl.GetAccountIDs()
	require.NoError(t, err)
	require.EqualValues(t, accountIDs, expectedAccountIDs)
}

func TestGetAccountIDsError(t *testing.T) {
	expectedAccountIDs := []string{"test_account_1"}

	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			ErrOnGetAccountIDs: true,
			AccountIDs:         expectedAccountIDs,
		}),
	})
	_ = ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})

	_ = ctrl.EnterPin("123123231")

	accountIDs, err := ctrl.GetAccountIDs()
	require.Error(t, err)
	require.Empty(t, accountIDs)
}

func TestSelectAccountID(t *testing.T) {
	selectedAccountID := "test_account_1"
	expectedAccountIDs := []string{selectedAccountID}

	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			AccountIDs: expectedAccountIDs,
		}),
	})
	_ = ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})

	_ = ctrl.EnterPin("123123231")

	err := ctrl.SelectAccount(selectedAccountID)
	require.NoError(t, err)
}

func TestSelectAccountIDInternalError(t *testing.T) {
	selectedAccountID := "test_account_1"
	expectedAccountIDs := []string{selectedAccountID}

	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			ErrOnSelectAccountID: true,
			AccountIDs:           expectedAccountIDs,
		}),
	})
	_ = ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})

	_ = ctrl.EnterPin("123123231")

	err := ctrl.SelectAccount(selectedAccountID)
	require.Error(t, err)
	require.EqualError(t, err, errorcode.FailedToSelectAccountID)
}

func TestSelectAccountID_DNE(t *testing.T) {
	selectedAccountID := "test_account_1"
	expectedAccountIDs := []string{selectedAccountID}

	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			AccountIDs: expectedAccountIDs,
		}),
	})
	_ = ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})

	_ = ctrl.EnterPin("123123231")

	err := ctrl.SelectAccount("12323")
	require.Error(t, err)
	require.EqualError(t, err, errorcode.NoMatchingAccountID)
}

func TestGetBalance(t *testing.T) {
	selectedAccountID := "test_account_1"
	expectedAccountIDs := []string{selectedAccountID}
	expectedBalanceAmt := 50

	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			AccountIDs:    expectedAccountIDs,
			GetBalanceAmt: expectedBalanceAmt,
		}),
	})
	_ = ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})

	_ = ctrl.EnterPin("123123231")

	err := ctrl.SelectAccount(selectedAccountID)
	require.NoError(t, err)

	balance, err := ctrl.GetBalance(selectedAccountID)
	require.NoError(t, err)
	require.Equal(t, expectedBalanceAmt, balance)
}

func TestGetBalanceError(t *testing.T) {
	selectedAccountID := "test_account_1"
	expectedAccountIDs := []string{selectedAccountID}
	expectedBalanceAmt := 50

	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			AccountIDs:      expectedAccountIDs,
			ErrOnGetBalance: true,
			GetBalanceAmt:   expectedBalanceAmt,
		}),
	})
	_ = ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})

	_ = ctrl.EnterPin("123123231")

	err := ctrl.SelectAccount(selectedAccountID)
	require.NoError(t, err)

	balance, err := ctrl.GetBalance(selectedAccountID)
	require.Error(t, err)
	require.Equal(t, -1, balance)
}

func TestMakeDeposit(t *testing.T) {
	selectedAccountID := "test_account_1"
	expectedAccountIDs := []string{selectedAccountID}

	expectedBalanceAmt := 50
	depositAmount := 30
	expectedBalanceAmtAfterDeposit := 80

	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			AccountIDs:          expectedAccountIDs,
			GetBalanceAmt:       expectedBalanceAmt,
			BalanceAfterDeposit: expectedBalanceAmtAfterDeposit,
		}),
	})
	_ = ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})

	_ = ctrl.EnterPin("123123231")

	err := ctrl.SelectAccount(selectedAccountID)
	require.NoError(t, err)

	balance, err := ctrl.GetBalance(selectedAccountID)
	require.NoError(t, err)
	require.Equal(t, expectedBalanceAmt, balance)

	newBalance, err := ctrl.MakeDeposit(selectedAccountID, depositAmount)
	require.NoError(t, err)
	require.Equal(t, expectedBalanceAmtAfterDeposit, newBalance)
}

func TestMakeDepositError(t *testing.T) {
	selectedAccountID := "test_account_1"
	expectedAccountIDs := []string{selectedAccountID}

	expectedBalanceAmt := 50
	depositAmount := 30
	expectedBalanceAmtAfterDeposit := 80

	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			AccountIDs:          expectedAccountIDs,
			GetBalanceAmt:       expectedBalanceAmt,
			ErrOnMakeDeposit:    true,
			BalanceAfterDeposit: expectedBalanceAmtAfterDeposit,
		}),
	})
	_ = ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})

	_ = ctrl.EnterPin("123123231")

	err := ctrl.SelectAccount(selectedAccountID)
	require.NoError(t, err)

	balance, err := ctrl.GetBalance(selectedAccountID)
	require.NoError(t, err)
	require.Equal(t, expectedBalanceAmt, balance)

	newBalance, err := ctrl.MakeDeposit(selectedAccountID, depositAmount)
	require.Error(t, err)
	require.Equal(t, -1, newBalance)
}

func TestMakeWithdrawal(t *testing.T) {
	selectedAccountID := "test_account_1"
	expectedAccountIDs := []string{selectedAccountID}

	expectedBalanceAmt := 50
	withdrawAmount := 30
	expectedBalanceAmtAfterWithdrawl := 20

	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			AccountIDs:           expectedAccountIDs,
			GetBalanceAmt:        expectedBalanceAmt,
			BalanceAfterWithdraw: expectedBalanceAmtAfterWithdrawl,
		}),
	})
	_ = ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})

	_ = ctrl.EnterPin("123123231")

	err := ctrl.SelectAccount(selectedAccountID)
	require.NoError(t, err)

	newBalance, err := ctrl.MakeWithdrawl(selectedAccountID, withdrawAmount)
	require.NoError(t, err)
	require.Equal(t, expectedBalanceAmtAfterWithdrawl, newBalance)
}

func TestMakeWithdrawInternalError(t *testing.T) {
	selectedAccountID := "test_account_1"
	expectedAccountIDs := []string{selectedAccountID}

	expectedBalanceAmt := 50
	withdrawAmount := 30

	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			AccountIDs:    expectedAccountIDs,
			GetBalanceAmt: expectedBalanceAmt,
			ErrOnWithdraw: true,
		}),
	})
	_ = ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})

	_ = ctrl.EnterPin("123123231")

	err := ctrl.SelectAccount(selectedAccountID)
	require.NoError(t, err)

	newBalance, err := ctrl.MakeWithdrawl(selectedAccountID, withdrawAmount)
	require.Error(t, err)
	require.Equal(t, -1, newBalance)
}

func TestMakeWithdraw_overdraw(t *testing.T) {
	selectedAccountID := "test_account_1"
	expectedAccountIDs := []string{selectedAccountID}

	expectedBalanceAmt := 50
	withdrawAmount := 80

	ctrl := NewAtmController(Options{
		cardSvc: testutil.NewDummyCardSvc(testutil.DummyCardTestOptions{}),
		accountSvc: testutil.NewDummyAccountSvc(testutil.DummyAcctTestOptions{
			AccountIDs:    expectedAccountIDs,
			GetBalanceAmt: expectedBalanceAmt,
		}),
	})
	_ = ctrl.InsertCard(model.Card{
		HolderName: "test user",
		Number:     "1234",
	})

	_ = ctrl.EnterPin("123123231")

	err := ctrl.SelectAccount(selectedAccountID)
	require.NoError(t, err)

	newBalance, err := ctrl.MakeWithdrawl(selectedAccountID, withdrawAmount)
	require.EqualError(t, err, errorcode.IsOverdraw)
	require.Equal(t, -1, newBalance)
}
