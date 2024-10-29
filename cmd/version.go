package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show application version",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", rootCmd.Name(), Version)
	},
}

func init() {
	// help flag
	versionCmd.Flags().BoolP("help", "h", false, "Help for command version")
}
