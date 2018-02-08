//
//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package version

import (
	"fmt"

	triton "github.com/joyent/triton-go"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:          "version",
	Short:        "print triton cli version",
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Version: %s\n", triton.UserAgent())
		return nil
	},
}
