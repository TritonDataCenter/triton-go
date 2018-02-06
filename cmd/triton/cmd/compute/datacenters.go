package compute

import (
	"context"

	"github.com/joyent/triton-go/cmd/internal/api"
	"github.com/joyent/triton-go/cmd/internal/console_writer"
	"github.com/joyent/triton-go/compute"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var DatacentersCommand = &cobra.Command{
	Use:   "datacenters",
	Short: "datacenters interaction with triton",
}

var ListDataCentersCommand = &cobra.Command{
	Use:          "list datacenters",
	Short:        "lists datacenters associated with triton account",
	SilenceUsage: true,
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

		datacenters, err := client.Datacenters().List(context.Background(), &compute.ListDataCentersInput{})
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

		table.SetHeader([]string{"name", "url"})

		var numDcs uint
		for _, dc := range datacenters {
			table.Append([]string{dc.Name, dc.URL})
			numDcs++
		}

		table.Render()

		return nil
	},
}
