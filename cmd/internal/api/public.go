//
//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package api

import (
	"github.com/joyent/triton-go"
	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/identity"
	"github.com/joyent/triton-go/network"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type TritonClientConfig struct {
	config *triton.ClientConfig
}

func (t *TritonClientConfig) GetAccountClient() (*account.AccountClient, error) {
	accountClient, err := account.NewClient(t.config)
	if err != nil {
		return nil, errors.Wrap(err, "Error Creating Triton Account Client")
	}

	return accountClient, nil
}

func (t *TritonClientConfig) GetComputeClient() (*compute.ComputeClient, error) {
	computeClient, err := compute.NewClient(t.config)
	if err != nil {
		return nil, errors.Wrap(err, "Error Creating Triton Compute Client")
	}
	return computeClient, nil
}

func (t *TritonClientConfig) GetIdentityClient() (*identity.IdentityClient, error) {
	identityClient, err := identity.NewClient(t.config)
	if err != nil {
		return nil, errors.Wrap(err, "Error Creating Triton Identity Client")
	}
	return identityClient, nil
}

func (t *TritonClientConfig) GetNetworkClient() (*network.NetworkClient, error) {
	networkClient, err := network.NewClient(t.config)
	if err != nil {
		return nil, errors.Wrap(err, "Error Creating Triton Netowkr Client")
	}
	return networkClient, nil
}

func InitConfig() (*TritonClientConfig, error) {
	viper.AutomaticEnv()

	var signer authentication.Signer
	var err error

	signer, err = authentication.NewSSHAgentSigner(authentication.SSHAgentSignerInput{
		KeyID:       viper.GetString("SDC_KEY_ID"),
		AccountName: viper.GetString("SDC_ACCOUNT"),
	})
	if err != nil {
		log.Fatal().Str("func", "initConfig").Msg("Error Creating SSH Agent Signer")
		return nil, err
	}

	config := &triton.ClientConfig{
		TritonURL:   viper.GetString("SDC_URL"),
		AccountName: viper.GetString("SDC_ACCOUNT"),
		Signers:     []authentication.Signer{signer},
	}

	return &TritonClientConfig{
		config: config,
	}, nil
}
