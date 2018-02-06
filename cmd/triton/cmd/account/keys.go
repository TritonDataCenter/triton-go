package account

import (
	"strings"

	"github.com/joyent/triton-go/cmd/internal/console_writer"
	"github.com/spf13/cobra"
)

var KeysCommand = &cobra.Command{
	Use:   "keys",
	Short: "key interaction with triton",
}

var ListKeysCommand = &cobra.Command{
	Use:          "list",
	Short:        "lists keys associated with triton account",
	SilenceUsage: true,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {

		cons := console_writer.GetTerminal()
		cons.Write([]byte(strings.Repeat("some crap\n", 1000)))

		//tritonClientConfig, err := api.InitConfig()
		//if err != nil {
		//	return err
		//}
		//
		//client, err := tritonClientConfig.GetAccountClient()
		//if err != nil {
		//	return err
		//}
		//
		//keys, err := client.Keys().List(context.Background(), &account.ListKeysInput{})
		//if err != nil {
		//	return err
		//}
		//
		//table := tablewriter.NewWriter(cons)
		//table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		//table.SetHeaderLine(false)
		//table.SetAutoFormatHeaders(true)
		//
		//table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
		//table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		//table.SetCenterSeparator("")
		//table.SetColumnSeparator("")
		//table.SetRowSeparator("")
		//
		//table.SetHeader([]string{"name", "fingerprint"})
		//
		//var numKeys uint
		//for _, key := range keys {
		//	table.Append([]string{key.Name, key.Fingerprint})
		//	numKeys++
		//}
		//
		//table.Render()

		return nil

	},
}
