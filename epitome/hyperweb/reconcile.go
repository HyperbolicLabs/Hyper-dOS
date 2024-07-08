package hyperweb

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func reconcile(
	clientset kubernetes.Clientset,
	dynamicClient dynamic.DynamicClient,
	gatewayUrl string,
	token string,
) error {
	if !secretExists(clientset, hyperwebNamespace, "operator-oauth") {
		// if it does not, query the gateway for oauth credentials using our token
		logrus.Infof("operator-oauth secret does not exist in namespace: %v", hyperwebNamespace)

		response, err := handshake(gatewayUrl, token)
		if err != nil {
			logrus.Errorf("failed to handshake with gateway: %v", err)
			return err
		}

		mustCreateOperatorOAuthSecret(
			clientset,
			hyperwebNamespace,
			"operator-oauth",
			response.ClientID,
			response.ClientSecret,
		)

		err = installClusterNameConfigMap(clientset, response.ClusterName)
		if err != nil {
			logrus.Errorf("failed to save cluster name in configmap: %v", err)
			return err
		}
	}

	name, err := GetClusterName(clientset)
	if err != nil {
		logrus.Errorf("failed to get cluster name: %v", err)
		return err
	}

	if IsInstalled(dynamicClient) {
		if isRegistered(clientset) {
			logrus.Infof("hyperweb application is installed and registered, nothing to do")
		} else {
			response, err := register(
				gatewayUrl,
				token,
				*name,
			)
			if err != nil {
				return err
			}
			if response.Success {
				logrus.Infof("registered cluster %s with gateway", *name)
			} else {
				return fmt.Errorf("failed to register cluster %s with gateway", *name)
			}
		}
	} else {
		logrus.Infof("hyperweb application is not installed - installing now")

		err = InstallHyperWeb(dynamicClient, *name)
		if err != nil {
			logrus.Errorf("failed to install hyperweb application: %v", err)
			return err
		}
	}

	return nil
}
