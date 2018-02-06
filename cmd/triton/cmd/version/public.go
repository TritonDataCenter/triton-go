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
