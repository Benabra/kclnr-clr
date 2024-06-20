package cmd

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func isDNSError(err error) bool {
	_, ok := err.(*net.DNSError)
	return ok || strings.Contains(err.Error(), "no such host")
}

func contextExists(contextName string) bool {
	kubeconfigPath := filepath.Join(homeDir(), ".kube", "config")
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		fmt.Printf("Failed to load kubeconfig: %v\n", err)
		return false
	}
	_, exists := config.Contexts[contextName]
	return exists
}

func checkContextStatus(contextName, kubeconfigPath string) (string, error) {
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}, configOverrides)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return "", err
	}

	_, err = clientset.ServerVersion()
	if err != nil {
		return "Unreachable", err
	}
	return "Reachable", nil
}
