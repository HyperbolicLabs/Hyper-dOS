package helper

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetSecret(
	clientset kubernetes.Interface,
	namespace string,
	name string,
) (*v1.Secret, error) {
	return clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
