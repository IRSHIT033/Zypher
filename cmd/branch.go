package cmd

import (
	"fmt"

	"github.com/IRSHIT033/zypher/pkg/zypher"
	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List, create, or delete branches",
	Long: `Manage branches in your Zypher repository.
Without arguments, lists all branches. With a branch name, creates a new branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// List branches
			if err := zypher.ListBranches(); err != nil {
				fmt.Printf("Error listing branches: %v\n", err)
				return
			}
		} else {
			// Create new branch
			branchName := args[0]
			if err := zypher.CreateNewBranch(branchName); err != nil {
				fmt.Printf("Error creating branch: %v\n", err)
				return
			}
		}
	},
}

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Switch branches",
	Long:  `Switch to another branch in your Zypher repository.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		if err := zypher.CheckoutBranch(branchName); err != nil {
			fmt.Printf("Error switching branch: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
	rootCmd.AddCommand(checkoutCmd)
}
