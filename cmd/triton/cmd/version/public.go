//
//  Copyright 2020 Joyent, Inc. All rights reserved.
//  Copyright 2025 MNX Cloud, Inc.
//  Copyright 2026 Edgecast Cloud LLC.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package version

import (
	"fmt"

	triton "github.com/TritonDataCenter/triton-go/v2"
	"github.com/TritonDataCenter/triton-go/v2/cmd/internal/command"
	"github.com/sean-/conswriter"
	"github.com/spf13/cobra"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:          "version",
		Short:        "print triton cli version",
		SilenceUsage: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			cons := conswriter.GetTerminal()
			fmt.Fprintf(cons, "Version: %s\n", triton.UserAgent())
			bi := triton.BuildInfo()
			if bi.Revision != "" {
				fmt.Fprintf(cons, "Revision: %s\n", bi.Revision)
			}
			if bi.Time != "" {
				fmt.Fprintf(cons, "Build Time: %s\n", bi.Time)
			}
			if bi.Modified {
				fmt.Fprintf(cons, "Modified: true\n")
			}
			return nil
		},
	},
	Setup: func(parent *command.Command) error {
		return nil
	},
}
