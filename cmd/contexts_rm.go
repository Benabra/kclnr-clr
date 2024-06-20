package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// removeContextCmd represents the rm command
var removeContextCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a kubeconfig context",
	Long:  `Removes a specified context or all unreachable contexts from the kubeconfig file.`,
	Run: func(cmd *cobra.Command, args []string) {
		contextName, _ := cmd.Flags().GetString("name")
		removeAllUnreachable, _ := cmd.Flags().GetBool("all-unreachable")

		if removeAllUnreachable {
			if confirm("Are you sure you want to remove all unreachable contexts? [y/N]: ") {
				if err := removeAllUnreachableContexts(); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("All unreachable contexts removed successfully.")
				}
			} else {
				fmt.Println("Operation cancelled.")
			}
		} else {
			if !contextExists(contextName) {
				fmt.Printf("Context %s does not exist.\n", contextName)
				return
			}

			if confirm(fmt.Sprintf("Are you sure you want to remove the context '%s'? [y/N]: ", contextName)) {
				if err := removeKubeconfigContext(contextName); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("Context %s removed successfully.\n", contextName)
				}
			} else {
				fmt.Println("Operation cancelled.")
			}
		}
	},
}

func init() {
	contextsCmd.AddCommand(removeContextCmd)
	removeContextCmd.Flags().StringP("name", "n", "", "Name of the context to remove")
	removeContextCmd.Flags().BoolP("all-unreachable", "a", false, "Remove all unreachable contexts")
}

func removeKubeconfigContext(contextName string) error {
	kubeconfigPath := filepath.Join(homeDir(), ".kube", "config")
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %v", err)
	}

	delete(config.Contexts, contextName)
	if config.CurrentContext == contextName {
		config.CurrentContext = ""
	}

	if err := clientcmd.WriteToFile(*config, kubeconfigPath); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %v", err)
	}

	return nil
}

func removeAllUnreachableContexts() error {
	kubeconfigPath := filepath.Join(homeDir(), ".kube", "config")
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %v", err)
	}

	unreachableContexts := []string{}
	for contextName := range config.Contexts {
		status, err := checkContextStatus(contextName, kubeconfigPath)
		if err != nil || status == "Unreachable" {
			fmt.Printf("%s: %s\n", contextName, status)
			unreachableContexts = append(unreachableContexts, contextName)
		}
	}

	for _, contextName := range unreachableContexts {
		delete(config.Contexts, contextName)
		if config.CurrentContext == contextName {
			config.CurrentContext = ""
		}
	}

	if err := clientcmd.WriteToFile(*config, kubeconfigPath); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %v", err)
	}

	return nil
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)
	return strings.ToLower(response) == "y"
}
