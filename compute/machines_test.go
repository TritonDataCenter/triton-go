package compute

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func getAnyMachineID(t *testing.T, c *Compute) (string, error) {
	machines, err := c.Machines().ListMachines(
		context.Background(),
		&ListMachinesInput{},
	)
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

func TestAccMachine_GetMachine(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{
			&testutils.StepAPICall{
				StateBagKey: "machine",
				CallFunc: func(client *Compute) (interface{}, error) {
					machineID, err := getAnyMachineID(t, client)
					if err != nil {
						return nil, err
					}

					return client.Machines().GetMachine(
						context.Background(),
						&GetMachineInput{
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
				CallFunc: func(client *Compute) (interface{}, error) {
					machineID, err := getAnyMachineID(t, client)
					if err != nil {
						return nil, err
					}

					return client.Machines().ListMachineTags(
						context.Background(),
						&testutils.ListMachineTagsInput{
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
