//
//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package cmd

import (
	"os"

	"github.com/joyent/triton-go/cmd/internal/command"
	"github.com/joyent/triton-go/cmd/internal/config"
	"github.com/joyent/triton-go/cmd/internal/console_writer"
	"github.com/joyent/triton-go/cmd/internal/logger"
	"github.com/joyent/triton-go/cmd/triton/cmd/autocompletion"
	"github.com/joyent/triton-go/cmd/triton/cmd/compute"
	"github.com/joyent/triton-go/cmd/triton/cmd/docs"
	"github.com/joyent/triton-go/cmd/triton/cmd/man"
	isatty "github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var subCommands = []*command.Command{
	compute.InstancesCommand,
	man.Cmd,
	autocompletion.Cmd,
	docs.Cmd,
}

var rootCmd = &command.Command{
	Cobra: &cobra.Command{
		Use:   "triton",
		Short: "cli for interacting with triton",
	},
	Setup: func(parent *command.Command) error {
		{
			const (
				key         = config.KeyUsePager
				longName    = "use-pager"
				shortName   = "P"
				description = "Use a pager to read the output (defaults to $PAGER, less(1), or more(1))"
			)
			var defaultValue bool
			if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
				defaultValue = true
			}

			flags := parent.Cobra.PersistentFlags()
			flags.BoolP(longName, shortName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = config.KeyManPageDirectory
				longName     = "dir"
				defaultValue = ""
				description  = "the directory to write the man pages"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		{
			const (
				key          = config.KeyMarkdownDirectory
				longName     = "docs-dir"
				defaultValue = ""
				description  = "the directory to write the markdown documentation pages"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		{
			const (
				key          = config.KeyAutoCompletionTarget
				longName     = "completionfile"
				defaultValue = "/etc/bash_completion.d/triton.sh"
				description  = "autocompletion file"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		{
			const (
				key          = config.KeyAutoCompletionType
				longName     = "type"
				defaultValue = "bash"
				description  = "autocompletion type (currently only bash supported)"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		return nil
	},
}

func Execute() error {

	rootCmd.Setup(rootCmd)

	console_writer.UsePager(viper.GetBool(config.KeyUsePager))

	if err := logger.Config(); err != nil {
		return err
	}

	for _, cmd := range subCommands {
		rootCmd.Cobra.AddCommand(cmd.Cobra)
		cmd.Setup(cmd)
	}

	if err := rootCmd.Cobra.Execute(); err != nil {
		//log.Error().Err(err).Msg("unable to run")
		return err
	}

	return nil
}

func generateDocumentation(cmd *cobra.Command) error {
	//err := doc.GenMarkdownTree(cmd, "./docs/md")
	//if err != nil {
	//	return err
	//}

	return nil
}
