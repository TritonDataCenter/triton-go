//
//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package man

import (
	"fmt"
	"os"

	"github.com/joyent/triton-go"
	"github.com/joyent/triton-go/cmd/internal/command"
	"github.com/joyent/triton-go/cmd/internal/config"
	"github.com/joyent/triton-go/cmd/internal/console_writer"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:   "man",
		Short: "Generates and installs triton cli man pages",
		Long: `This command automatically generates up-to-date man pages of Triton CLI
command-line interface.  By default, it creates the man page files
in the "docs/man" directory under the current directory.`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cons := console_writer.GetTerminal()

			header := &doc.GenManHeader{
				Manual:  "Triton",
				Section: "3",
				Source:  fmt.Sprintf("Triton %s", triton.Version),
			}

			location := viper.GetString(config.KeyManPageDirectory)
			if location == "" {
				location = "docs/man"
			}
			if _, err := os.Stat(location); os.IsNotExist(err) {
				os.Mkdir(location, 0777)
			}

			cmd.Root().DisableAutoGenTag = true
			cons.Write([]byte(fmt.Sprintf("Generating manpages to %s", location)))

			err := doc.GenManTree(cmd.Root(), header, location)
			if err != nil {
				return err
			}

			cons.Write([]byte("\nManpage generation complete"))

			return nil
		},
	},
	Setup: func(parent *command.Command) error {
		return nil
	},
}
