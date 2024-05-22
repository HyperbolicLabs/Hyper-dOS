package hyperweb

import (
	"k8s.io/client-go/kubernetes"
)

const hyperwebNamespace = "hyperweb"

func SecretExists(
	clientset kubernetes.Clientset,
	namespace string,
	name string,
) bool {
	_, err := GetSecret(clientset, namespace, name)
	return err == nil
}
