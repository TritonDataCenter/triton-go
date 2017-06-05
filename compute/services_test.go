package compute

import (
	"context"
	"fmt"
	"testing"
)

func TestAccServicesList(t *testing.T) {
	const stateKey = "services"

	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{
			&testutils.StepAPICall{
				StateBagKey: stateKey,
				CallFunc: func(client *Compute) (interface{}, error) {
					return client.Services().ListServices(
						context.Background(),
						&ListServicesInput{})
				},
			},
			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
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
