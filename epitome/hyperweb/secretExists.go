package hyperweb

import (
	"context"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const hyperwebNamespace = "hyperweb"

func Initialize(clientset kubernetes.Clientset) {
	// check if operator-oauth.hyperweb.secrets.cluster.local exists
	if !SecretExists(clientset, hyperwebNamespace, "operator-oauth") {
		logrus.Fatalf("operator-oauth secret not found in namespace: %v", hyperwebNamespace)
	}
}

func SecretExists(
	clientset kubernetes.Clientset,
	namespace string,
	name string,
) bool {
	_, err := GetSecret(clientset, namespace, name)
	return err == nil
}

func GetSecret(
	clientset kubernetes.Clientset,
	namespace string,
	name string,
) (*v1.Secret, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return secret, nil
}
