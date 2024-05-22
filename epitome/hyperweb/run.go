package hyperweb

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

func Run(clientset kubernetes.Clientset, gatewayUrl string, token string) {
	// check if operator-oauth.hyperweb.secrets.cluster.local exists
	if !SecretExists(clientset, hyperwebNamespace, "operator-oauth") {
		// if it does not, query the gateway for oauth credentials using our token

		clientId, clientSecret, err := handshake(gatewayUrl, token)
		if err != nil {
			logrus.Warnf("failed to handshake with gateway: %v", err)
			return
		} else {
			mustCreateOperatorOAuthSecret(
				clientset,
				hyperwebNamespace,
				"operator-oauth",
				*clientId,
				*clientSecret,
			)
		}
	}
}
