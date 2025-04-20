package cmd

import (
	"fmt"

	"github.com/IRSHIT033/zypher/pkg/zypher"
	"github.com/spf13/cobra"
)

var revertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Revert to a specific commit",
	Long:  `Revert your working directory to the state of a specific commit.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hash := args[0]
		if err := zypher.RevertToCommit(hash); err != nil {
			fmt.Printf("Error reverting to commit: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(revertCmd)
}
