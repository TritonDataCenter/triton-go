package account

import "github.com/joyent/triton-go/client"

type AccountService struct {
	client *client.Client
}

// Accounts returns a Compute client used for accessing functions pertaining to
// Account functionality in the Triton API.
func (c *AccountService) Accounts() *AccountsClient {
	return &AccountsClient{c}
}

// Config returns a c used for accessing functions pertaining
// to Config functionality in the Triton API.
func (c *AccountService) Config() *ConfigClient {
	return &ConfigClient{c}
}

// Keys returns a Compute client used for accessing functions pertaining to SSH
// key functionality in the Triton API.
func (c *AccountService) Keys() *KeysClient {
	return &KeysClient{c}
}
