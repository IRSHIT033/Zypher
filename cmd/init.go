package cmd

import (
	"fmt"

	"github.com/IRSHIT033/zypher/pkg/zypher"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Zypher repository",
	Long:  `Initialize a new Zypher repository in the current directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := zypher.InitRepository(); err != nil {
			fmt.Printf("Error initializing repository: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
