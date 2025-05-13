package cluster

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GenerateClientsets(
	kubeconfig *string,
) (kubernetes.Interface, *dynamic.DynamicClient, error) {
	if kubeconfig != nil {
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create clientconfig: %v", err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create clientset: %v", err)
		}

		dynamicClient, err := dynamic.NewForConfig(config)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create dynamic clientset: %v", err)
		}

		return clientset, dynamicClient, nil
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create in-cluster config: %v", err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create dynamic clientset: %v", err)
	}

	return clientset, dynamicClient, nil
}
