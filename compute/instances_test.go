package compute

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/testutils"
)

func getAnyInstanceID(t *testing.T, client *compute.ComputeClient) (string, error) {
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

func TestAccInstances_Get(t *testing.T) {
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

					instanceID, err := getAnyInstanceID(t, c)
					if err != nil {
						return nil, err
					}

					ctx := context.Background()
					input := &compute.GetInstanceInput{
						ID: instanceID,
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
// returned from getAnyInstanceID will have at least one tag.
func TestAccInstances_ListTags(t *testing.T) {
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

					instanceID, err := getAnyInstanceID(t, c)
					if err != nil {
						return nil, err
					}

					ctx := context.Background()
					input := &compute.ListTagsInput{
						ID: instanceID,
					}
					return c.Instances().ListTags(ctx, input)
				},
			},

			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
					tagsRaw, found := state.GetOk("instances")
					if !found {
						return fmt.Errorf("State key %q not found", "instances")
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

func TestAccInstances_UpdateMetadata(t *testing.T) {
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

					instanceID, err := getAnyInstanceID(t, c)
					if err != nil {
						return nil, err
					}

					ctx := context.Background()
					input := &compute.UpdateMetadataInput{
						ID: instanceID,
						Metadata: map[string]string{
							"tester": os.Getenv("USER"),
						},
					}
					return c.Instances().UpdateMetadata(ctx, input)
				},
			},

			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
					mdataRaw, found := state.GetOk("instances")
					if !found {
						return fmt.Errorf("State key %q not found", "instances")
					}

					mdata := mdataRaw.(map[string]string)
					if len(mdata) == 0 {
						return errors.New("Expected metadata on machine")
					}

					if mdata["tester"] != os.Getenv("USER") {
						return errors.New("Expected test metadata to equal environ $USER")
					}
					return nil
				},
			},
		},
	})
}

func TestAccInstances_ListMetadata(t *testing.T) {
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

					instanceID, err := getAnyInstanceID(t, c)
					if err != nil {
						return nil, err
					}

					ctx := context.Background()
					input := &compute.ListMetadataInput{
						ID: instanceID,
					}
					return c.Instances().ListMetadata(ctx, input)
				},
			},

			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
					mdataRaw, found := state.GetOk("instances")
					if !found {
						return fmt.Errorf("State key %q not found", "instances")
					}

					mdata := mdataRaw.(map[string]string)
					if len(mdata) == 0 {
						return errors.New("Expected metadata on machine")
					}

					if mdata["root_authorized_keys"] == "" {
						return errors.New("Expected test metadata to have key")
					}
					return nil
				},
			},
		},
	})
}

func TestAccInstances_GetMetadata(t *testing.T) {
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

					instanceID, err := getAnyInstanceID(t, c)
					if err != nil {
						return nil, err
					}

					ctx := context.Background()
					input := &compute.UpdateMetadataInput{
						ID: instanceID,
						Metadata: map[string]string{
							"testkey": os.Getenv("USER"),
						},
					}
					_, err = c.Instances().UpdateMetadata(ctx, input)
					if err != nil {
						return nil, err
					}

					ctx2 := context.Background()
					input2 := &compute.GetMetadataInput{
						ID:  instanceID,
						Key: "testkey",
					}
					return c.Instances().GetMetadata(ctx2, input2)
				},
			},

			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
					mdataValue := state.Get("instances")
					retValue := fmt.Sprintf("\"%s\"", os.Getenv("USER"))
					if mdataValue != retValue {
						return errors.New("Expected test metadata to equal environ \"$USER\"")
					}
					return nil
				},
			},
		},
	})
}
