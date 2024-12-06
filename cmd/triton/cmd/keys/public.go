//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
// Copyright 2024 MNX Cloud, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package keys

import (
	"github.com/TritonDataCenter/triton-go/cmd/internal/command"
	"github.com/TritonDataCenter/triton-go/cmd/internal/config"
	"github.com/TritonDataCenter/triton-go/cmd/triton/cmd/keys/create"
	"github.com/TritonDataCenter/triton-go/cmd/triton/cmd/keys/delete"
	"github.com/TritonDataCenter/triton-go/cmd/triton/cmd/keys/get"
	"github.com/TritonDataCenter/triton-go/cmd/triton/cmd/keys/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:     "keys",
		Aliases: []string{"key"},
		Short:   "List and manage Triton SSH Keys.",
	},

	Setup: func(parent *command.Command) error {

		cmds := []*command.Command{
			list.Cmd,
			get.Cmd,
			keyDelete.Cmd,
			create.Cmd,
		}

		for _, cmd := range cmds {
			cmd.Setup(cmd)
			parent.Cobra.AddCommand(cmd.Cobra)
		}

		{
			const (
				key          = config.KeySSHKeyFingerprint
				longName     = "fingerprint"
				defaultValue = ""
				description  = "SSH Key Fingerprint"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		{
			const (
				key          = config.KeySSHKeyName
				longName     = "keyname"
				defaultValue = ""
				description  = "SSH Key Name"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		return nil
	},
}
