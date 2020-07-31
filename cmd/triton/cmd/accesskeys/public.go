//
// Copyright 2020 Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package accesskeys

import (
	"github.com/joyent/triton-go/v2/cmd/internal/command"
	"github.com/joyent/triton-go/v2/cmd/internal/config"
	"github.com/joyent/triton-go/v2/cmd/triton/cmd/accesskeys/create"
	accessKeyDelete "github.com/joyent/triton-go/v2/cmd/triton/cmd/accesskeys/delete"
	"github.com/joyent/triton-go/v2/cmd/triton/cmd/accesskeys/get"
	"github.com/joyent/triton-go/v2/cmd/triton/cmd/accesskeys/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:     "accesskeys",
		Aliases: []string{"accesskey"},
		Short:   "List and manage Triton Access Keys.",
	},

	Setup: func(parent *command.Command) error {

		cmds := []*command.Command{
			list.Cmd,
			get.Cmd,
			accessKeyDelete.Cmd,
			create.Cmd,
		}

		for _, cmd := range cmds {
			cmd.Setup(cmd)
			parent.Cobra.AddCommand(cmd.Cobra)
		}

		{
			const (
				key          = config.KeyAccessKeyID
				longName     = "accesskeyid"
				defaultValue = ""
				description  = "Access Key Identifier"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		return nil
	},
}
