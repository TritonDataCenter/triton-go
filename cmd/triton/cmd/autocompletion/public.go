//
//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package autocompletion

import (
	"errors"
	"fmt"

	"github.com/joyent/triton-go/cmd/internal/command"
	"github.com/joyent/triton-go/cmd/internal/config"
	"github.com/joyent/triton-go/cmd/internal/console_writer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:   "autocomplete",
		Short: "Generates shell autocompletion file for Triton",
		Long: `Generates a shell autocompletion script for Triton.

NOTE: The current version supports Bash only.
      This should work for *nix systems with Bash installed.

By default, the file is written directly to /etc/bash_completion.d
for convenience, and the command may need superuser rights, e.g.:

	$ sudo triton bashcompletion

Add ` + "`--completionfile=/path/to/file`" + ` flag to set alternative
file-path and name.

Logout and in again to reload the completion scripts,
or just source them in directly:

	$ . /etc/bash_completion`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cons := console_writer.GetTerminal()

			if viper.GetString(config.KeyAutoCompletionType) != "bash" {
				return errors.New("Only Bash is support for Autocomplete Type")
			}

			target := viper.GetString(config.KeyAutoCompletionTarget)
			err := cmd.Root().GenBashCompletionFile(target)
			if err != nil {
				return err
			}

			cons.Write([]byte(fmt.Sprintf("Bash completion file saved to %s", target)))

			return nil
		},
	},
	Setup: func(parent *command.Command) error {
		return nil
	},
}
