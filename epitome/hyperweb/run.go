package hyperweb

import (
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

var clusterName *string

func Runloop(
	clientset kubernetes.Clientset,
	dynamicClient dynamic.DynamicClient,
	gatewayUrl string,
	token string,
) {
	for {
		// run once every 30 seconds
		reconcile(clientset, dynamicClient, gatewayUrl, token)
		time.Sleep(30 * time.Second)
	}
}

func reconcile(
	clientset kubernetes.Clientset,
	dynamicClient dynamic.DynamicClient,
	gatewayUrl string,
	token string,
) {
	if !secretExists(clientset, hyperwebNamespace, "operator-oauth") {
		// if it does not, query the gateway for oauth credentials using our token
		logrus.Infof("operator-oauth secret does not exist in namespace: %v", hyperwebNamespace)

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
		)

		clusterName = &response.ClusterName

		err = InstallCM(dynamicClient, *clusterName)
		if err != nil {
			logrus.Fatalf("failed to save cluster name in configmap: %v", err)
		}

		InstallHyperWeb(dynamicClient, *clusterName)
	}

	if IsInstalled(dynamicClient) {
		logrus.Infof("hyperweb application is installed, nothing to do")
		return
	} else {
		logrus.Infof("hyperweb application is not installed")
	}

	// if !applicationExists(clientset, "argocd", "hyperweb") {
	// 	mustCreateApplication(clientset, "argocd", "hyperweb", hyperwebNamespace)
	// }
}
