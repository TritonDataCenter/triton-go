//
//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package compute

//var DatacentersCommand = &cobra.Command{
//	Use:   "datacenters",
//	Short: "datacenters interaction with triton",
//}
//
//var ListDataCentersCommand = &cobra.Command{
//	Args:         cobra.NoArgs,
//	Use:          "list datacenters",
//	Short:        "lists datacenters associated with triton account",
//	SilenceUsage: true,
//	PreRunE: func(cmd *cobra.Command, args []string) error {
//		return nil
//	},
//	RunE: func(cmd *cobra.Command, args []string) error {
//		cons := console_writer.GetTerminal()
//
//		tritonClientConfig, err := api.InitConfig()
//		if err != nil {
//			return err
//		}
//
//		client, err := tritonClientConfig.GetComputeClient()
//		if err != nil {
//			return err
//		}
//
//		datacenters, err := client.Datacenters().List(context.Background(), &compute.ListDataCentersInput{})
//		if err != nil {
//			return err
//		}
//
//		table := tablewriter.NewWriter(cons)
//		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
//		table.SetHeaderLine(false)
//		table.SetAutoFormatHeaders(true)
//
//		table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
//		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
//		table.SetCenterSeparator("")
//		table.SetColumnSeparator("")
//		table.SetRowSeparator("")
//
//		table.SetHeader([]string{"name", "url"})
//
//		var numDcs uint
//		for _, dc := range datacenters {
//			table.Append([]string{dc.Name, dc.URL})
//			numDcs++
//		}
//
//		table.Render()
//
//		return nil
//	},
//}
