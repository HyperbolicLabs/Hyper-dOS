package helper

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfigMap(
	clientset kubernetes.Interface,
	namespace string,
	name string,
) (*v1.ConfigMap, error) {
	return clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
