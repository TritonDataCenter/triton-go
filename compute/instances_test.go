package compute

import (
	"context"
	"errors"
	"fmt"
	"testing"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/testutils"
)

func getAnyMachineID(t *testing.T, client *compute.ComputeClient) (string, error) {
	ctx := context.Background()
	input := &compute.ListInstancesInput{}
	instances, err := client.Instances().List(ctx, input)
	if err != nil {
		return "", err
	}

	for _, m := range instances {
		if len(m.ID) > 0 {
			return m.ID, nil
		}
	}

	t.Skip()
	return "", errors.New("no machines configured")
}

func TestAccMachine_Get(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{

			&testutils.StepClient{
				StateBagKey: "instances",
				CallFunc: func(config *triton.ClientConfig) (interface{}, error) {
					return compute.NewClient(config)
				},
			},

			&testutils.StepAPICall{
				StateBagKey: "instances",
				CallFunc: func(client interface{}) (interface{}, error) {
					c := client.(*compute.ComputeClient)

					machineID, err := getAnyMachineID(t, c)
					if err != nil {
						return nil, err
					}

					ctx := context.Background()
					input := &compute.GetInstancesInput{
						ID: machineID,
					}
					return c.Instances().Get(ctx, input)
				},
			},

			&testutils.StepAssertSet{
				StateBagKey: "instances",
				Keys:        []string{"ID", "Name", "Type", "Tags"},
			},
		},
	})
}

// FIXME(seanc@): TestAccMachine_ListMachineTags assumes that any machine ID
// returned from getAnyMachineID will have at least one tag.
func TestAccMachine_ListTags(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{

			&testutils.StepClient{
				StateBagKey: "instances",
				CallFunc: func(config *triton.ClientConfig) (interface{}, error) {
					return compute.NewClient(config)
				},
			},

			&testutils.StepAPICall{
				StateBagKey: "instances",
				CallFunc: func(client interface{}) (interface{}, error) {
					c := client.(*compute.ComputeClient)

					machineID, err := getAnyMachineID(t, c)
					if err != nil {
						return nil, err
					}

					ctx := context.Background()
					input := &compute.ListTagsInput{
						ID: machineID,
					}
					return c.Instances().ListTags(ctx, input)
				},
			},

			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
					tagsRaw, found := state.GetOk("instances")
					if !found {
						return fmt.Errorf("State key %q not found", "machines")
					}

					tags := tagsRaw.(map[string]interface{})
					if len(tags) == 0 {
						return errors.New("Expected at least one tag on machine")
					}
					return nil
				},
			},
		},
	})
}
