//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package account

import (
	"github.com/joyent/triton-go/client"
	"github.com/joyent/triton-go/triton"
)

type AccountClient struct {
	Client *client.Client
}

func newAccountClient(client *client.Client) *AccountClient {
	return &AccountClient{
		Client: client,
	}
}

// NewClient returns a new client for working with Account endpoints and
// resources within CloudAPI
func NewClient(config *triton.ClientConfig) (*AccountClient, error) {
	// TODO: Utilize config interface within the function itself
	client, err := client.New(config.TritonURL, config.MantaURL, config.AccountName, config.Signers...)
	if err != nil {
		return nil, err
	}
	return newAccountClient(client), nil
}

// Config returns a c used for accessing functions pertaining
// to Config functionality in the Triton API.
func (c *AccountClient) Config() *ConfigClient {
	return &ConfigClient{c.Client}
}

// Keys returns a Compute client used for accessing functions pertaining to SSH
// key functionality in the Triton API.
func (c *AccountClient) Keys() *KeysClient {
	return &KeysClient{c.Client}
}
