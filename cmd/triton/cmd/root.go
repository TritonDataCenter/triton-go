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
	"github.com/joyent/triton-go/cmd/internal/logger"
	"github.com/joyent/triton-go/cmd/triton/cmd/docs"
	"github.com/joyent/triton-go/cmd/triton/cmd/instances"
	"github.com/joyent/triton-go/cmd/triton/cmd/shell"
	isatty "github.com/mattn/go-isatty"
	"github.com/sean-/conswriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var subCommands = []*command.Command{
	instances.Cmd,
	docs.Cmd,
	shell.Cmd,
}

var rootCmd = &command.Command{
	Cobra: &cobra.Command{
		Use:   "triton",
		Short: "cli for interacting with triton",
		Long: `CLI for Joyent Triton. In order to interact with Triton, 
you will need to have the following environment variables set:

* SDC_KEY_ID or TRITON_KEY_ID
* SDC_ACCOUNT or TRITON_ACCOUNT
* SDC_URL or TRITON_URL

This CLI doesn't have support for passing key material. It will try and
locate your SSH KEY from ssh-agent.`,
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
				key          = config.KeyLogLevel
				longOpt      = "log-level"
				shortOpt     = "l"
				defaultValue = "INFO"
				description  = "Change the log level being sent to stdout"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.StringP(longOpt, shortOpt, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longOpt))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key         = config.KeyLogFormat
				longOpt     = "log-format"
				shortOpt    = "F"
				description = `Specify the log format ("auto", "zerolog", or "human")`
			)
			defaultValue := logger.FormatAuto.String()

			flags := parent.Cobra.PersistentFlags()
			flags.StringP(longOpt, shortOpt, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longOpt))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key         = config.KeyLogTermColor
				longOpt     = "use-color"
				shortOpt    = ""
				description = "Use ASCII colors"
			)
			defaultValue := false
			if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
				defaultValue = true
			}

			flags := parent.Cobra.PersistentFlags()
			flags.BoolP(longOpt, shortOpt, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longOpt))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = config.KeyUseUTC
				longName     = "utc"
				shortName    = "Z"
				defaultValue = false
				description  = "Display times in UTC"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.BoolP(longName, shortName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}

func Execute() error {

	rootCmd.Setup(rootCmd)

	conswriter.UsePager(viper.GetBool(config.KeyUsePager))

	if err := logger.Setup(); err != nil {
		return err
	}

	for _, cmd := range subCommands {
		rootCmd.Cobra.AddCommand(cmd.Cobra)
		cmd.Setup(cmd)
	}

	if err := rootCmd.Cobra.Execute(); err != nil {
		return err
	}

	return nil
}
