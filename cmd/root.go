package cmd

import (
	"fmt"
	"os"

	"github.com/IRSHIT033/zypher/pkg/zypher"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zypher",
	Short: "Zypher is a simple version control system",
	Long: `Zypher is a simple version control system inspired by Git.
It provides basic version control functionality for tracking changes in your projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show logo and help when no subcommand is provided
		zypher.PrintLogo()
		fmt.Println()
		fmt.Println("Usage: zypher <command> [options]")
		fmt.Println("\nCommands:")
		fmt.Println("  init      Initialize a new repository")
		fmt.Println("  status    Show the status of your working directory")
		fmt.Println("  commit    Create a new commit")
		fmt.Println("  revert    Revert to a specific commit")
		fmt.Println("  log       Show commit history")
		fmt.Println("  branch    List, create, or delete branches")
		fmt.Println("  checkout  Switch branches")
		fmt.Println("\nUse 'zypher help <command>' for more information about a command.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
