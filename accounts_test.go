package triton

import (
	"context"
	"testing"
)

func TestAccAccount_Get(t *testing.T) {
	AccTest(t, TestCase{
		Steps: []Step{
			&StepAPICall{
				StateBagKey: "account",
				CallFunc: func(client *Client) (interface{}, error) {
					return client.Accounts().GetAccount(
						context.Background(), &GetAccountInput{})
				},
			},
			&StepAssertSet{
				StateBagKey: "account",
				Keys:        []string{"ID", "Login", "Email"},
			},
		},
	})
}
