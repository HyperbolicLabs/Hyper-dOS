package shtest

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func InitFakeCluster() kubernetes.Interface {
	clientset := fake.NewSimpleClientset()

	for _, ns := range testNamespaces {
		clientset.CoreV1().Namespaces().Create(
			context.TODO(),
			&corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: ns,
				},
			},
			metav1.CreateOptions{},
		)
	}

	return clientset
}

var testNamespaces = []string{
	"ping",
	"hyperdos",
	"instance",
	"hyperweb",
	"argocd",
	"tailscale-operator",
	"hyperbolic",
	"calico-system",
}
