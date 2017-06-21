package compute

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func getAnyMachineID(t *testing.T, c *InstancesClient) (string, error) {
	machines, err := c.Instances().List(context.Background(), &ListInstancesInput{})
	if err != nil {
		return "", err
	}

	for _, m := range machines {
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
			&testutils.StepAPICall{
				StateBagKey: "machine",
				CallFunc: func(client *InstancesClient) (interface{}, error) {
					machineID, err := getAnyMachineID(t, client)
					if err != nil {
						return nil, err
					}

					return client.Instances().Get(
						context.Background(),
						&GetInstanceInput{
							ID: machineID,
						})
				},
			},
			&testutils.StepAssertSet{
				StateBagKey: "machine",
				Keys:        []string{"ID", "Name", "Type", "Tags"},
			},
		},
	})
}

// FIXME(seanc@): TestAccMachine_ListMachineTags assumes that any machine ID
// returned from getAnyMachineID will have at least one tag.
func TestAccMachine_ListMachineTags(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{
			&testutils.StepAPICall{
				StateBagKey: "machine",
				CallFunc: func(client *InstancesClient) (interface{}, error) {
					machineID, err := getAnyMachineID(t, client)
					if err != nil {
						return nil, err
					}

					return client.Instances().ListTags(
						context.Background(),
						&ListTagsInput{
							ID: machineID,
						})
				},
			},
			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
					tagsRaw, found := state.GetOk("machine")
					if !found {
						return fmt.Errorf("State key %q not found", "machines")
					}

					tags := tagsRaw.(map[string]string)
					if len(tags) == 0 {
						return errors.New("Expected at least one tag on machine")
					}
					return nil
				},
			},
		},
	})
}
