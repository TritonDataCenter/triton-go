package triton

import (
	"context"
	"fmt"
	"testing"
)

// Note that this is specific to Joyent Public Cloud and will not pass on
// private installations of Triton.
func TestAccNetworks_List(t *testing.T) {
	AccTest(t, TestCase{
		Steps: []Step{
			&StepAPICall{
				StateBagKey: "networks",
				CallFunc: func(client *Client) (interface{}, error) {
					return client.Networks().ListNetworks(
						context.Background(),
						&ListNetworksInput{})
				},
			},
			&StepAssertFunc{
				AssertFunc: func(state TritonStateBag) error {
					dcs, ok := state.GetOk("networks")
					if !ok {
						return fmt.Errorf("State key %q not found", "datacenters")
					}

					toFind := []string{"Joyent-SDC-Private", "Joyent-SDC-Public"}
					for _, dcName := range toFind {
						found := false
						for _, dc := range dcs.([]*Network) {
							if dc.Name == dcName {
								found = true
								if dc.Id == "" {
									return fmt.Errorf("%q has no ID", dc.Name)
								}
							}
						}
						if !found {
							return fmt.Errorf("Did not find Network %q", dcName)
						}
					}

					return nil
				},
			},
		},
	})
}
