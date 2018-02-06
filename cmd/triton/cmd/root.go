package cmd

import (
	"github.com/joyent/triton-go/cmd/internal/config"
	"github.com/joyent/triton-go/cmd/internal/console_writer"
	"github.com/joyent/triton-go/cmd/internal/logger"
	"github.com/joyent/triton-go/cmd/triton/cmd/account"
	"github.com/joyent/triton-go/cmd/triton/cmd/compute"
	"github.com/joyent/triton-go/cmd/triton/cmd/identity"
	"github.com/joyent/triton-go/cmd/triton/cmd/network"
	"github.com/joyent/triton-go/cmd/triton/cmd/version"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "triton",
	Short: "cli for interacting with triton",
}

func Execute() error {
	console_writer.UsePager(viper.GetBool(config.KeyUsePager))

	if err := logger.Config(); err != nil {
		return err
	}

	rootCmd.AddCommand(account.AccountCommand)
	rootCmd.AddCommand(identity.IdentityCommand)
	rootCmd.AddCommand(network.NetworkCommand)
	rootCmd.AddCommand(version.Cmd)

	account.SetUpCommands()
	compute.SetUpCommands(rootCmd)
	identity.SetUpCommands()
	network.SetUpCommands()

	initRootFlags()

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("unable to run")
		return err
	}

	return nil
}

func initRootFlags() {
	{
		const (
			key          = config.KeyUsePager
			longName     = "use-pager"
			shortName    = "P"
			defaultValue = false
			description  = "Use a pager to read the output (defaults to $PAGER, less(1), or more(1))"
		)

		rootCmd.PersistentFlags().BoolP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, rootCmd.PersistentFlags().Lookup(longName))
		viper.BindEnv(key, "TRITON_USE_PAGER")
		viper.SetDefault(key, defaultValue)
	}
}
