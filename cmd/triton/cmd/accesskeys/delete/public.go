//
// Copyright 2020 Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package accesskeydelete

import (
	"errors"
	"fmt"

	"github.com/joyent/triton-go/v2/cmd/agent/account"
	cfg "github.com/joyent/triton-go/v2/cmd/config"
	"github.com/joyent/triton-go/v2/cmd/internal/command"
	"github.com/sean-/conswriter"
	"github.com/spf13/cobra"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Args:         cobra.NoArgs,
		Use:          "delete",
		Short:        "delete Triton Access Key",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cfg.GetAccessKeyID() == "" {
				return errors.New("`accesskeyid` must be specified")
			}

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

			accesskey, err := a.DeleteAccessKey()
			if err != nil {
				return err
			}

			cons.Write([]byte(fmt.Sprintf("Deleted access key %q", accesskey.AccessKeyID)))

			return nil
		},
	},
	Setup: func(parent *command.Command) error {
		return nil
	},
}
