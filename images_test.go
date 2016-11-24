package triton

import (
	"fmt"
	"testing"

	"github.com/abdullin/seq"
)

func TestAccImages_List(t *testing.T) {
	const stateKey = "images"
	const image1Id = "95f6c9a6-a2bd-11e2-b753-dbf2651bf890"
	const image2Id = "70e3ae72-96b6-11e6-9056-9737fd4d0764"
	AccTest(t, TestCase{
		Steps: []Step{
			&StepAPICall{
				StateBagKey: stateKey,
				CallFunc: func(client *Client) (interface{}, error) {
					return client.Images().ListImages(&ListImagesInput{})
				},
			},
			&StepAssertFunc{
				AssertFunc: func(state TritonStateBag) error {
					images, ok := state.GetOk(stateKey)
					if !ok {
						return fmt.Errorf("State key %q not found", stateKey)
					}

					toFind := []string{image1Id, image2Id}
					for _, imageID := range toFind {
						found := false
						for _, image := range images.([]*Image) {
							if image.ID == imageID {
								found = true
								state.Put(imageID, image)
							}
						}
						if !found {
							return fmt.Errorf("Did not find Image %q", imageID)
						}
					}

					return nil
				},
			},
			&StepAssert{
				StateBagKey: image1Id,
				Assertions: seq.Map{
					"name":                    "ws2012std",
					"owner":                   "9dce1460-0c4c-4417-ab8b-25ca478c5a78",
					"requirements.min_memory": 3840,
					"requirements.min_ram":    3840,
				},
			},
			&StepAssert{
				StateBagKey: image2Id,
				Assertions: seq.Map{
					"name":       "base-64",
					"owner":      "9dce1460-0c4c-4417-ab8b-25ca478c5a78",
					"tags.role":  "os",
					"tags.group": "base-64",
				},
			},
		},
	})
}
