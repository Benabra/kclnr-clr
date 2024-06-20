package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// listNodesCmd represents the list command for nodes
var listNodesCmd = &cobra.Command{
	Use:   "list",
	Short: "List Kubernetes nodes and check their status",
	Long:  `Lists all available Kubernetes nodes and checks their status.`,
	Run: func(cmd *cobra.Command, args []string) {
		contextName, _ := cmd.Flags().GetString("context")
		if err := listNodesAndCheckStatus(contextName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	nodesCmd.AddCommand(listNodesCmd)
	listNodesCmd.Flags().StringP("context", "c", "", "Name of the kubeconfig context")
}

func listNodesAndCheckStatus(contextName string) error {
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
			if isDNSError(err) {
				data = append(data, []string{ctxName, "-", "DNS resolution error"})
			} else {
				data = append(data, []string{ctxName, "-", fmt.Sprintf("Error: %v", err)})
			}
			continue
		}

		clientset, err := kubernetes.NewForConfig(restConfig)
		if err != nil {
			data = append(data, []string{ctxName, "-", fmt.Sprintf("Error: %v", err)})
			continue
		}

		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			if isDNSError(err) {
				data = append(data, []string{ctxName, "-", "DNS resolution error"})
			} else {
				data = append(data, []string{ctxName, "-", fmt.Sprintf("Error: %v", err)})
			}
			continue
		}

		for _, node := range nodes.Items {
			status := "Unknown"
			for _, condition := range node.Status.Conditions {
				if condition.Type == "Ready" {
					if condition.Status == "True" {
						status = color.GreenString("Ready")
					} else {
						status = color.RedString("NotReady")
					}
					break
				}
			}
			data = append(data, []string{ctxName, node.Name, status})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Context", "Node Name", "Status"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return nil
}
