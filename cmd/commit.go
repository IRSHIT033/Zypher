package cmd

import (
	"fmt"

	"github.com/IRSHIT033/zypher/pkg/zypher"
	"github.com/spf13/cobra"
)

var message string

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a new commit",
	Long:  `Create a new commit with the specified message.`,
	Run: func(cmd *cobra.Command, args []string) {
		if message == "" {
			fmt.Println("Error: Commit message is required")
			return
		}
		if err := zypher.CreateCommit(message); err != nil {
			fmt.Printf("Error creating commit: %v\n", err)
			return
		}
	},
}

func init() {
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message")
	rootCmd.AddCommand(commitCmd)
}
