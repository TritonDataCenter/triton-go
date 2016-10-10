package triton

import (
	"fmt"
	"testing"

	"github.com/abdullin/seq"
)

// Note that this is specific to Joyent Public Cloud and will not pass on
// private installations of Triton.
func TestAccDataCenters_Get(t *testing.T) {
	const dataCenterName = "us-east-1"
	const dataCenterURL = "https://us-east-1.api.joyentcloud.com"

	AccTest(t, TestCase{
		Steps: []Step{
			&StepAPICall{
				StateBagKey: "datacenter",
				CallFunc: func(client *Client) (interface{}, error) {
					return client.Datacenters().GetDataCenter(&GetDataCenterInput{
						Name: dataCenterName,
					})
				},
			},
			&StepAssert{
				StateBagKey: "datacenter",
				Assertions: seq.Map{
					"name": dataCenterName,
					"url":  dataCenterURL,
				},
			},
		},
	})
}

// Note that this is specific to Joyent Public Cloud and will not pass on
// private installations of Triton.
func TestAccDataCenters_List(t *testing.T) {
	AccTest(t, TestCase{
		Steps: []Step{
			&StepAPICall{
				StateBagKey: "datacenters",
				CallFunc: func(client *Client) (interface{}, error) {
					return client.Datacenters().ListDataCenters(&ListDataCentersInput{})
				},
			},
			&StepAssertFunc{
				AssertFunc: func(state TritonStateBag) error {
					dcs, ok := state.GetOk("datacenters")
					if !ok {
						return fmt.Errorf("State key %q not found", "datacenters")
					}

					toFind := []string{"us-east-1", "eu-ams-1"}
					for _, dcName := range toFind {
						found := false
						for _, dc := range dcs.([]*DataCenter) {
							if dc.Name == dcName {
								found = true
								if dc.URL == "" {
									return fmt.Errorf("%q has no URL", dc.Name)
								}
							}
						}
						if !found {
							return fmt.Errorf("Did not find DC %q", dcName)
						}
					}

					return nil
				},
			},
		},
	})
}
