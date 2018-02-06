package compute

import (
	"context"

	"github.com/joyent/triton-go/cmd/internal/api"
	"github.com/joyent/triton-go/cmd/internal/console_writer"
	"github.com/joyent/triton-go/compute"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var ImagesCommand = &cobra.Command{
	Use:   "images",
	Short: "image interaction with triton",
}

var ListImagesCommand = &cobra.Command{
	Use:          "list images",
	Short:        "lists images associated with triton account",
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

		images, err := client.Images().List(context.Background(), &compute.ListImagesInput{})
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

		table.SetHeader([]string{"id", "name", "description", "version"})

		var numImages uint
		for _, image := range images {
			table.Append([]string{image.ID, image.Name, image.Description, image.Version})
			numImages++
		}

		table.Render()

		return nil
	},
}
