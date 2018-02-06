package compute

import (
	"context"

	"github.com/joyent/triton-go/cmd/internal/api"
	"github.com/joyent/triton-go/cmd/internal/console_writer"
	"github.com/joyent/triton-go/compute"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var InstancesCommand = &cobra.Command{
	Use:        "instance",
	Aliases:    []string{"instances"},
	SuggestFor: []string{"machines"},
	Short:      "instance interaction with triton",
}

var ListInstancesCommand = &cobra.Command{
	Use:          "list istances",
	Short:        "lists instances associated with triton account",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cons := console_writer.GetTerminal()

		tritonClientConfig, err := api.InitConfig()
		if err != nil {
			return err
		}

		client, err := tritonClientConfig.GetComputeClient()
		if err != nil {
			return err
		}

		instances, err := client.Instances().List(context.Background(), &compute.ListInstancesInput{})
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(cons)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeaderLine(false)
		table.SetAutoFormatHeaders(true)

		table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")

		table.SetHeader([]string{"id", "name", "image", "package"})

		var numInstances uint
		for _, instance := range instances {
			table.Append([]string{instance.ID, instance.Name, instance.Image, instance.Package})
			numInstances++
		}

		table.Render()

		return nil
	},
}
