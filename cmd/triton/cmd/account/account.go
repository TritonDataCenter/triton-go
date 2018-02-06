package account

import (
	"github.com/spf13/cobra"
)

var AccountCommand = &cobra.Command{
	Use:   "account",
	Short: "account interaction with triton",
}

func SetUpCommands() {
	AccountCommand.AddCommand(KeysCommand)
	KeysCommand.AddCommand(ListKeysCommand)
}
