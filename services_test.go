package triton

import (
	"fmt"
	"testing"
)

func TestAccServicesList(t *testing.T) {
	const stateKey = "services"

	AccTest(t, TestCase{
		Steps: []Step{
			&StepAPICall{
				StateBagKey: stateKey,
				CallFunc: func(client *Client) (interface{}, error) {
					return client.Services().ListServices(&ListServicesInput{})
				},
			},
			&StepAssertFunc{
				AssertFunc: func(state TritonStateBag) error {
					services, ok := state.GetOk(stateKey)
					if !ok {
						return fmt.Errorf("State key %q not found", stateKey)
					}

					toFind := []string{"docker"}
					for _, serviceName := range toFind {
						found := false
						for _, service := range services.([]*Service) {
							if service.Name == serviceName {
								found = true
								if service.Endpoint == "" {
									return fmt.Errorf("%q has no Endpoint", service.Name)
								}
							}
						}
						if !found {
							return fmt.Errorf("Did not find Service %q", serviceName)
						}
					}

					return nil
				},
			},
		},
	})
}
