package account

import (
	"context"
	"testing"
)

func TestAccConfig_Get(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{
			&testutils.StepAPICall{
				StateBagKey: "config",
				CallFunc: func(client *AccountClient) (interface{}, error) {
					return client.Config().Get(context.Background(), &GetConfigInput{})
				},
			},
			&testutils.StepAssertSet{
				StateBagKey: "config",
				Keys:        []string{"DefaultNetwork"},
			},
		},
	})
}
