package triton

import (
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/client"
	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/storage"
)

type TritonClient struct {
	*client.Client
	Compute *compute.Compute
	Storage *storage.Storage
}

type ClientConfig struct {
	*client.Config
	Endpoint    string
	AccountName string
	Signers     []authentication.Signer
}

// TODO: Work configuration providers into the mix for pulling variables out of
// the env or `node-triton` profile.
func NewClient(config *ClientConfig) (*TritonClient, error) {
	// TODO: Utilize config interface within the function itself
	client, error := client.New(config.Endpoint, config.AccountName, config.Signers...)
	if error != nil {
		return nil, error
	}
	compute := &compute.Compute{client}
	storage := &storage.Storage{client}
	return &TritonClient{client, compute, storage}, nil
}
