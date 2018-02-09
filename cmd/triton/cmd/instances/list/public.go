//
//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package list

import (
	"github.com/joyent/triton-go/cmd/agent/compute"
	cfg "github.com/joyent/triton-go/cmd/config"
	"github.com/joyent/triton-go/cmd/internal/command"
	"github.com/joyent/triton-go/cmd/internal/config"
	"github.com/olekukonko/tablewriter"
	"github.com/sean-/conswriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Args:         cobra.NoArgs,
		Use:          "list",
		Short:        "list triton instances",
		Aliases:      []string{"ls"},
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cons := conswriter.GetTerminal()

			c, err := cfg.New()
			if err != nil {
				return err
			}

			a, err := compute.NewGetComputeClient(c)
			if err != nil {
				return err
			}

			instances, err := a.GetInstanceList()
			if err != nil {
				return err
			}

			images, err := a.GetImagesList()
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

			table.SetHeader([]string{"SHORTID", "NAME", "IMG", "STATE", "AGE"})

			var numInstances uint
			for _, instance := range instances {
				table.Append([]string{string(instance.ID[:8]), instance.Name, a.FormatImageName(images, instance.Image), instance.State, "--"})
				numInstances++
			}

			table.Render()

			return nil
		},
	},
	Setup: func(parent *command.Command) error {
		{
			const (
				key          = config.KeyInstanceNamePrefix
				longName     = "name-prefix"
				defaultValue = ""
				description  = "Instance Name Prefix"
			)

			flags := parent.Cobra.Flags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		{
			const (
				key          = config.KeyInstanceState
				longName     = "state"
				defaultValue = ""
				description  = "Instance state (e.g. running)"
			)

			flags := parent.Cobra.Flags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		{
			const (
				key          = config.KeyInstanceBrand
				longName     = "brand"
				defaultValue = ""
				description  = "Instance brand (e.g. lx, kvm)"
			)

			flags := parent.Cobra.Flags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		{
			const (
				key         = config.KeyInstanceSearchTag
				longName    = "tags"
				description = "Filter instances based on tag. This option can be used multiple times."
			)

			flags := parent.Cobra.Flags()
			flags.StringSlice(longName, nil, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		return nil
	},
}
