//
//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  Copyright 2024 MNX Cloud, Inc.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package network

import (
	"github.com/TritonDataCenter/triton-go/cmd/config"
	"github.com/TritonDataCenter/triton-go/network"
	"github.com/pkg/errors"
)

func NewGetNetworkClient(cfg *config.TritonClientConfig) (*network.NetworkClient, error) {
	networkClient, err := network.NewClient(cfg.Config)
	if err != nil {
		return nil, errors.Wrap(err, "Error Creating Triton Netowkr Client")
	}
	return networkClient, nil
}
