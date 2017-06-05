package compute

import (
	"context"
	"testing"
)

func TestAccConfig_Get(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{
			&testutils.StepAPICall{
				StateBagKey: "config",
				CallFunc: func(client *Compute) (interface{}, error) {
					return client.Config().GetConfig(
						context.Background(),
						&GetConfigInput{})
				},
			},
			&testutils.StepAssertSet{
				StateBagKey: "config",
				Keys:        []string{"DefaultNetwork"},
			},
		},
	})
}
