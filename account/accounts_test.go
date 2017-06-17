package account

import (
	"context"
	"testing"

	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/testutils"
)

func TestAccAccount_Get(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{
			&testutils.StepAPICall{
				StateBagKey: "account",
				CallFunc: func(client *AccountService) (interface{}, error) {
					return client.Accounts().GetAccount(
						context.Background(), &account.GetAccountInput{})
				},
			},
			&testutils.StepAssertSet{
				StateBagKey: "account",
				Keys:        []string{"ID", "Login", "Email"},
			},
		},
	})
}
