package jungle

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (a *agent) reconcile() error {
	if !secretExists(a.clientset, a.cfg.HyperwebNamespace, "operator-oauth") {
		// if it does not, query the gateway for oauth credentials using our token
		logrus.Infof("operator-oauth secret does not exist in namespace: %v", a.cfg.HyperwebNamespace)

		response, err := a.handshake(http.DefaultClient)
		if err != nil {
			logrus.Errorf("failed to handshake with gateway: %v", err)
			return err
		}

		a.mustCreateTailscaleOperatorOAuthSecret(
			response.ClientID,
			response.ClientSecret,
		)

		err = a.installClusterNameConfigMap(response.ClusterName)
		if err != nil {
			logrus.Errorf("failed to save cluster name in configmap: %v", err)
			return err
		}
	}

	name, err := a.getClusterName()
	if err != nil {
		logrus.Errorf("failed to get cluster name: %v", err)
		return err
	}

	if a.HyperwebIsInstalled() {
		if a.isRegistered() {
			logrus.Infof("hyperweb application is installed and registered, nothing to do")
		} else {
			response, err := a.register(
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

		err = a.installHyperWeb(*name)
		if err != nil {
			logrus.Errorf("failed to install hyperweb application: %v", err)
			return err
		}
	}
	return nil
}
