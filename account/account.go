package account

import "github.com/joyent/triton-go/client"

type Account struct {
	client *client.Client
}

// Accounts returns a Compute client used for accessing functions pertaining to
// Account functionality in the Triton API.
func (c *Compute) Accounts() *AccountsClient {
	return &AccountsClient{c}
}

// Config returns a c used for accessing functions pertaining
// to Config functionality in the Triton API.
func (c *Compute) Config() *ConfigClient {
	return &ConfigClient{c}
}

// Keys returns a Compute client used for accessing functions pertaining to SSH
// key functionality in the Triton API.
func (c *Compute) Keys() *KeysClient {
	return &KeysClient{c}
}
