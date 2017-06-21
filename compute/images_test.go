package compute

import (
	"fmt"
	"testing"

	"context"
	"time"

	"github.com/abdullin/seq"
	"github.com/joyent/triton-go/testutils"
)

func TestAccImagesList(t *testing.T) {
	const stateKey = "images"
	const image1Id = "95f6c9a6-a2bd-11e2-b753-dbf2651bf890"
	const image2Id = "70e3ae72-96b6-11e6-9056-9737fd4d0764"
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{
			&testutils.StepAPICall{
				StateBagKey: stateKey,
				CallFunc: func(client *ImagesClient) (interface{}, error) {
					return client.Images().List(
						context.Background(), &ListInput{})
				},
			},
			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
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
			&testutils.StepAssert{
				StateBagKey: image1Id,
				Assertions: seq.Map{
					"name":                    "ws2012std",
					"owner":                   "9dce1460-0c4c-4417-ab8b-25ca478c5a78",
					"requirements.min_memory": 3840,
					"requirements.min_ram":    3840,
				},
			},
			&testutils.StepAssert{
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

func TestAccImagesGet(t *testing.T) {
	const stateKey = "image"
	const imageId = "95f6c9a6-a2bd-11e2-b753-dbf2651bf890"
	publishedAt, err := time.Parse(time.RFC3339, "2013-04-11T21:07:38Z")
	if err != nil {
		t.Fatal("Reference time does not parse as RFC3339")
	}
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{
			&testutils.StepAPICall{
				StateBagKey: stateKey,
				CallFunc: func(client *ImagesClient) (interface{}, error) {
					return client.Images().Get(
						context.Background(),
						&GetInput{
							ImageID: imageId,
						})
				},
			},
			&testutils.StepAssert{
				StateBagKey: stateKey,
				Assertions: seq.Map{
					"name":    "ws2012std",
					"version": "1.0.1",
					"os":      "windows",
					"requirements.min_memory": 3840,
					"requirements.min_ram":    3840,
					"type":                    "zvol",
					"description":             "Windows Server 2012 Standard 64-bit image.",
					"files[0].compression":    "gzip",
					"files[0].sha1":           "fe35a3b70f0a6f8e5252b05a35ee397d37d15185",
					"files[0].size":           3985823590,
					"tags.role":               "os",
					"published_at":            publishedAt,
					"owner":                   "9dce1460-0c4c-4417-ab8b-25ca478c5a78",
					"public":                  true,
					"state":                   "active",
				},
			},
		},
	})
}
