//
//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package get

import (
	"errors"
	"fmt"

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
		Use:          "get",
		Short:        "get a triton instance",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cfg.GetMachineID() == "" && cfg.GetMachineName() == "" {
				return errors.New("Either `id` or `name` must be specified")
			}

			if cfg.GetMachineID() != "" && cfg.GetMachineName() != "" {
				return errors.New("Only 1 of `id` or `name` must be specified")
			}

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

			instance, err := a.GetInstance()
			if err != nil {
				return err
			}

			table := tablewriter.NewWriter(cons)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetHeaderLine(false)
			table.SetAutoFormatHeaders(true)

			table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")

			table.SetHeader([]string{"------", "------"})

			table.Append([]string{"id", instance.ID})
			table.Append([]string{"name", instance.Name})
			table.Append([]string{"package", instance.Package})
			table.Append([]string{"image", instance.Image})
			table.Append([]string{"brand", instance.Brand})
			table.Append([]string{"firewall enabled", fmt.Sprintf("%t", instance.FirewallEnabled)})

			table.Render()

			return nil
		},
	},

	Setup: func(parent *command.Command) error {
		{
			const (
				key          = config.KeyInstanceName
				longName     = "name"
				shortName    = "n"
				defaultValue = ""
				description  = "Instance Name"
			)

			flags := parent.Cobra.Flags()
			flags.StringP(longName, shortName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		{
			const (
				key          = config.KeyInstanceID
				longName     = "id"
				defaultValue = ""
				description  = "Instance ID"
			)

			flags := parent.Cobra.Flags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		return nil
	},
}
