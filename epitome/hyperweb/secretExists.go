package hyperweb

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

const hyperwebNamespace = "hyperweb"

func secretExists(
	clientset kubernetes.Clientset,
	namespace string,
	name string,
) bool {
	_, err := GetSecret(clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			return false
		}
		logrus.Fatalf("unexpected error trying to get secret: %v", err)
	}
	return true
}
