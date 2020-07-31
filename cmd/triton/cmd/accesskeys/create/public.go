//
// Copyright 2020 Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package create

import (
	"fmt"

	"github.com/joyent/triton-go/v2/cmd/agent/account"
	cfg "github.com/joyent/triton-go/v2/cmd/config"
	"github.com/joyent/triton-go/v2/cmd/internal/command"
	"github.com/joyent/triton-go/v2/cmd/internal/config"
	"github.com/sean-/conswriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Args:         cobra.NoArgs,
		Use:          "create",
		Aliases:      []string{"add"},
		Short:        "create Triton Access Key",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cons := conswriter.GetTerminal()

			c, err := cfg.NewTritonConfig()
			if err != nil {
				return err
			}

			a, err := account.NewAccountClient(c)
			if err != nil {
				return err
			}

			accesskey, err := a.CreateAccessKey()
			if err != nil {
				return err
			}

			cons.Write([]byte(fmt.Sprintf("Created access key %q", accesskey.AccessKeyID)))

			return nil
		},
	},
	Setup: func(parent *command.Command) error {

		{
			const (
				key          = config.KeyAccessKeyID
				longName     = "accesskeyid"
				defaultValue = ""
				description  = "Access Key Identifier"
			)

			flags := parent.Cobra.Flags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		return nil
	},
}
