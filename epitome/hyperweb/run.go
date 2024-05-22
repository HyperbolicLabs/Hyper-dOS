package hyperweb

import (
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

var clusterName *string

func Runloop(clientset kubernetes.Clientset, gatewayUrl string, token string) {

	for {
		// run once every 30 seconds
		reconcile(clientset, gatewayUrl, token)
		time.Sleep(30 * time.Second)
	}
}

func reconcile(clientset kubernetes.Clientset, gatewayUrl string, token string) {
	if !secretExists(clientset, hyperwebNamespace, "operator-oauth") {
		// if it does not, query the gateway for oauth credentials using our token

		response, err := handshake(gatewayUrl, token)
		if err != nil {
			logrus.Warnf("failed to handshake with gateway: %v", err)
			return
		}

		mustCreateOperatorOAuthSecret(
			clientset,
			hyperwebNamespace,
			"operator-oauth",
			response.ClientID,
			response.ClientSecret,
			response.ClusterName,
		)

		clusterName = &response.ClusterName
	}

	logrus.Info("TODO: check if hyperweb is installed")

	// if !applicationExists(clientset, "argocd", "hyperweb") {
	// 	mustCreateApplication(clientset, "argocd", "hyperweb", hyperwebNamespace)
	// }
}
