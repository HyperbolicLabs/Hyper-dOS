package cluster

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func MustConnect(
	kubeconfig *string,
) (
	kubernetes.Interface,
	*dynamic.DynamicClient) {

	clientset, dynamicClient, err := GenerateClientsets(kubeconfig)
	if err != nil {
		panic(err)
	}

	return clientset, dynamicClient
}
