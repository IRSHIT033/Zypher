package cmd

import (
	"fmt"

	"github.com/IRSHIT033/zypher/pkg/zypher"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of your working directory",
	Long:  `Show the status of your working directory, including modified files and staged changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := zypher.ShowStatus(); err != nil {
			fmt.Printf("Error showing status: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
