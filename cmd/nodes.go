package cmd

import (
	"github.com/spf13/cobra"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Manage Kubernetes nodes",
	Long:  `Manage Kubernetes nodes, including listing and checking their status.`,
}

func init() {
	contextsCmd.AddCommand(nodesCmd)
}
