package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// listContextsCmd represents the list command
var listContextsCmd = &cobra.Command{
	Use:   "list",
	Short: "List kubeconfig contexts and check their status",
	Long:  `Lists all available kubeconfig contexts from the kubeconfig file and checks the status of each context.`,
	Run: func(cmd *cobra.Command, args []string) {
		contextName, _ := cmd.Flags().GetString("context")
		if err := listKubeconfigContextsAndCheckStatus(contextName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	contextsCmd.AddCommand(listContextsCmd)
	listContextsCmd.Flags().StringP("context", "c", "", "Name of the kubeconfig context to check")
}

func listKubeconfigContextsAndCheckStatus(contextName string) error {
	kubeconfigPath := filepath.Join(homeDir(), ".kube", "config")
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %v", err)
	}

	data := [][]string{}
	contextsToCheck := map[string]*clientcmdapi.Context{}
	if contextName != "" {
		ctx, exists := config.Contexts[contextName]
		if !exists {
			return fmt.Errorf("context %s does not exist", contextName)
		}
		contextsToCheck[contextName] = ctx
	} else {
		contextsToCheck = config.Contexts
	}

	for ctxName, _ := range contextsToCheck {
		configOverrides := &clientcmd.ConfigOverrides{CurrentContext: ctxName}
		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}, configOverrides)

		restConfig, err := clientConfig.ClientConfig()
		if err != nil {
			data = append(data, []string{ctxName, color.RedString("Unreachable (%v)", err)})
			continue
		}

		clientset, err := kubernetes.NewForConfig(restConfig)
		if err != nil {
			data = append(data, []string{ctxName, color.RedString("Unreachable (%v)", err)})
			continue
		}

		_, err = clientset.ServerVersion()
		if err != nil {
			data = append(data, []string{ctxName, color.RedString("Unreachable (%v)", err)})
		} else {
			data = append(data, []string{ctxName, color.GreenString("Reachable")})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Status"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil
}
