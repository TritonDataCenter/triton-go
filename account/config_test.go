package account

import (
	"context"
	"testing"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/testutils"
)

func TestAccConfig_Get(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{

			&testutils.StepClient{
				StateBagKey: "config",
				CallFunc: func(config *triton.ClientConfig) (interface{}, error) {
					return account.NewClient(config)
				},
			},

			&testutils.StepAPICall{
				StateBagKey: "config",
				CallFunc: func(client interface{}) (interface{}, error) {
					c := client.(*account.AccountClient)
					ctx := context.Background()
					input := &account.GetConfigInput{}
					return c.Config().Get(ctx, input)
				},
			},
			&testutils.StepAssertSet{
				StateBagKey: "config",
				Keys:        []string{"DefaultNetwork"},
			},
		},
	})
}
