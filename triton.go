package triton

import (
	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/client"
	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/identity"
	"github.com/joyent/triton-go/network"
	"github.com/joyent/triton-go/storage"
)

type TritonClient struct {
	*client.Client
	Account  *account.AccountService
	Compute  *compute.Compute
	Identity *identity.IdentityService
	Network  *network.NetworkService
	Storage  *storage.Storage
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
	account := &account.AccountService{client}
	compute := &compute.Compute{client}
	identity := &identity.IdentityService{client}
	network := &network.NetworkService{client}
	storage := &storage.Storage{client}
	triton := &TritonClient{client, account, compute, identity, network, storage}
	return triton, nil
}
