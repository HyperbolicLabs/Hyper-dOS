package cluster

import (
	argo "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func MustConnect(
	kubeconfig *string,
) (
	kubernetes.Interface,
	*dynamic.DynamicClient,
	argo.Interface,
) {

	clientset, dynamicClient, argoClient, err := GenerateClientsets(kubeconfig)
	if err != nil {
		panic(err)
	}

	return clientset, dynamicClient, argoClient
}
