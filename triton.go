package triton

import (
	"github.com/joyent/triton-go/client"
	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/storage"
)

type TritonClient struct {
	*client.Client
	Compute compute.Compute
	Storage storage.Storage
}

func NewClient(config *client.Config) (*TritonClient, error) {
	// TODO: Utilize config interface within the function itself
	client, error := client.New(config.endpoint, config.accountName, config.signers)
	if error != nil {
		return nil, error
	}
	compute := &Compute{client}
	storage := &Storage{client}
	return &TritonClient{client, compute, storage}
}
