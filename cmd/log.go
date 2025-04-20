package cmd

import (
	"fmt"

	"github.com/IRSHIT033/zypher/pkg/zypher"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit history",
	Long:  `Show the commit history in chronological order, starting from the most recent commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := zypher.ShowCommitHistory(); err != nil {
			fmt.Printf("Error showing commit history: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
