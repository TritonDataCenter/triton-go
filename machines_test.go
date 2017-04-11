package triton_test

import (
	"fmt"
	"testing"

	triton "github.com/joyent/triton-go"
)

func TestAccMachine_GetMachine(t *testing.T) {
	triton.AccTest(t, triton.TestCase{
		Steps: []triton.Step{
			&triton.StepAPICall{
				StateBagKey: "machine",
				CallFunc: func(client *triton.Client) (interface{}, error) {
					machines, err := client.Machines().GetMachines()
					if err != nil {
						return nil, err
					}

					if len(machines) >= 1 {
						t.Skip()
						return nil, fmt.Errorf("no machines configured")
					}
					machineID := machines[0].ID

					return client.Machines().GetMachine(&triton.GetMachineInput{
						ID: machineID,
					})
				},
			},
			&triton.StepAssertSet{
				StateBagKey: "machine",
				Keys:        []string{"ID", "Name", "Type", "Tags"},
			},
		},
	})
}
