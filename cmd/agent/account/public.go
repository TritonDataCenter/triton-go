//
//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package account

import (
	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/cmd/config"
	"github.com/pkg/errors"
)

func NewGetAccountClient(cfg *config.TritonClientConfig) (*account.AccountClient, error) {
	accountClient, err := account.NewClient(cfg.Config)
	if err != nil {
		return nil, errors.Wrap(err, "Error Creating Triton Account Client")
	}

	return accountClient, nil
}
