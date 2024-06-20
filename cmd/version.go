package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Set the version for the application
const version = "0.15.4"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of kclnr",
	Long:  `All software has versions. This is kclnr's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kclnr version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
