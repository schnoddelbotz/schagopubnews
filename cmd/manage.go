package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// manageCmd represents the manage command
var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Manage newspapers and users for SPN",
	Long:  `After deployment, create a newspaper and a user like this: ...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("manage called")
	},
}

func init() {
	rootCmd.AddCommand(manageCmd)
}
