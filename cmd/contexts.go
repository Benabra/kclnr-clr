package cmd

import (
	"github.com/spf13/cobra"
)

// contextsCmd represents the contexts command
var contextsCmd = &cobra.Command{
	Use:   "contexts",
	Short: "Manage kubeconfig contexts",
	Long:  `Manage kubeconfig contexts, including listing and removing contexts.`,
}

func init() {
	rootCmd.AddCommand(contextsCmd)
}
