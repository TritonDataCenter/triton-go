package compute

import (
	"github.com/spf13/cobra"
)

var ComputeCommand = &cobra.Command{
	Use:   "compute",
	Short: "compute interaction with triton",
}

func SetUpCommands() {
	ComputeCommand.AddCommand(InstancesCommand)
	InstancesCommand.AddCommand(ListInstancesCommand)

	ComputeCommand.AddCommand(ImagesCommand)
	ImagesCommand.AddCommand(ListImagesCommand)

	ComputeCommand.AddCommand(DatacentersCommand)
	DatacentersCommand.AddCommand(ListDataCentersCommand)
}
