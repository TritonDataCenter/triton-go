package compute

import (
	"github.com/joyent/triton-go/cmd/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SetUpCommands(rootCommand *cobra.Command) {
	initComputeFlags()

	rootCommand.AddCommand(InstancesCommand)
	InstancesCommand.AddCommand(ListInstancesCommand)
	InstancesCommand.AddCommand(GetInstanceCommand)
	InstancesCommand.AddCommand(CreateInstanceCommand)

	//rootCommand.AddCommand(ImagesCommand)
	//ImagesCommand.AddCommand(ListImagesCommand)

	//rootCommand.AddCommand(DatacentersCommand)
	//DatacentersCommand.AddCommand(ListDataCentersCommand)
}

func initComputeFlags() {
	{
		const (
			key          = config.KeyInstanceId
			longName     = "id"
			defaultValue = ""
			description  = "Instance id (defaults to '')"
		)

		GetInstanceCommand.Flags().String(longName, defaultValue, description)
		viper.BindPFlag(key, GetInstanceCommand.Flags().Lookup(longName))
	}

	{
		const (
			key          = config.KeyInstanceName
			longName     = "name"
			shortName    = "N"
			defaultValue = ""
			description  = "Instance name (defaults to '')"
		)

		GetInstanceCommand.Flags().StringP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, GetInstanceCommand.Flags().Lookup(longName))

		CreateInstanceCommand.Flags().StringP(longName, shortName, defaultValue, description)
		CreateInstanceCommand.MarkFlagRequired(longName)
		viper.BindPFlag(key, CreateInstanceCommand.Flags().Lookup(longName))
	}

	{
		const (
			key          = config.KeyPackageName
			longName     = "pkg-name"
			defaultValue = ""
			description  = "Package name (defaults to '')"
		)

		CreateInstanceCommand.Flags().String(longName, defaultValue, description)
		viper.BindPFlag(key, CreateInstanceCommand.Flags().Lookup(longName))
	}

	{
		const (
			key          = config.KeyPackageId
			longName     = "pkg-id"
			defaultValue = ""
			description  = "Package id (defaults to ''). This takes precedence over 'pkg-name'"
		)

		CreateInstanceCommand.Flags().String(longName, defaultValue, description)
		viper.BindPFlag(key, CreateInstanceCommand.Flags().Lookup(longName))
	}

	{
		const (
			key          = config.KeyImageId
			longName     = "img-id"
			defaultValue = ""
			description  = "Image id (defaults to ''). This takes precedence over 'img-name'"
		)

		CreateInstanceCommand.Flags().String(longName, defaultValue, description)
		viper.BindPFlag(key, CreateInstanceCommand.Flags().Lookup(longName))
	}

	{
		const (
			key          = config.KeyImageName
			longName     = "img-name"
			defaultValue = ""
			description  = "Image name (defaults to '')"
		)

		CreateInstanceCommand.Flags().String(longName, defaultValue, description)
		viper.BindPFlag(key, CreateInstanceCommand.Flags().Lookup(longName))
	}

	{
		const (
			key          = config.KeyInstanceWait
			longName     = "wait"
			shortName    = "w"
			defaultValue = false
			description  = "Wait for the creation to complete (defaults to false)"
		)

		CreateInstanceCommand.Flags().BoolP(longName, shortName, defaultValue, description)
		viper.BindPFlag(key, CreateInstanceCommand.Flags().Lookup(longName))

		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = config.KeyInstanceFirewall
			longName     = "firewall"
			defaultValue = false
			description  = "Enable Cloud Firewall on this instance (defaults to false)"
		)

		CreateInstanceCommand.Flags().Bool(longName, defaultValue, description)
		viper.BindPFlag(key, CreateInstanceCommand.Flags().Lookup(longName))

		viper.SetDefault(key, defaultValue)
	}
}
