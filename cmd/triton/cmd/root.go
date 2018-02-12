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
		Short: "Joyent Triton CLI and client (https://www.joyent.com/triton)",
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

		{
			const (
				key          = config.KeyAccount
				longName     = "account"
				shortName    = "a"
				defaultValue = ""
				description  = "Account (login name). If not specified, the environment variable TRITON_ACCOUNT or SDC_ACCOUNT will be used"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.StringP(longName, shortName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		{
			const (
				key          = config.KeyUrl
				longName     = "url"
				shortName    = "U"
				defaultValue = ""
				description  = "CloudAPI URL. If not specified, the environment variable TRITON_URL or SDC_URL will be used"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.StringP(longName, shortName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		{
			const (
				key          = config.KeySshKeyID
				longName     = "key-id"
				shortName    = "k"
				defaultValue = ""
				description  = "SSH key fingerprint. If not specified, the environment variable TRITON_KEY_ID or SDC_KEY_ID will be used"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.StringP(longName, shortName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
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
