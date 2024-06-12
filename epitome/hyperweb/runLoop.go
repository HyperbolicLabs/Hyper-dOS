package hyperweb

import (
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func RunLoop(
	clientset kubernetes.Clientset,
	dynamicClient dynamic.DynamicClient,
	gatewayUrl string,
	token string,
	interval time.Duration,
) {
	for {
		// run once every 30 seconds
		err := reconcile(clientset, dynamicClient, gatewayUrl, token)
		if err != nil {
			logrus.Fatalf("failed to reconcile: %v", err)
		}
		time.Sleep(interval)
	}
}
