package jungle

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (a *agent) reconcile() error {
	if !secretExists(a.clientset, a.cfg.HyperwebNamespace, "operator-oauth") {
		// if it does not, query the gateway for oauth credentials using our token
		a.logger.Info(
			"operator-oauth secret does not exist, will handshake",
			zap.String("namespace", a.cfg.HyperwebNamespace))

		response, err := a.handshake(http.DefaultClient)
		if err != nil {
			return fmt.Errorf("failed to handshake: %v", err)
		}

		err = a.createTailscaleOperatorOAuthSecret(
			response.ClientID,
			response.ClientSecret,
		)
		if err != nil {
			return fmt.Errorf("failed to create operator-oauth secret: %v", err)
		}

		err = a.installClusterNameConfigMap(response.ClusterName)
		if err != nil {
			return fmt.Errorf("failed to save cluster name in configmap: %v", err)
		}
	}

	name, err := a.getClusterName()
	if err != nil {
		return fmt.Errorf("failed to get cluster name: %v", err)
	}

	if a.HyperwebIsInstalled() {
		if a.isRegistered() {
			a.logger.Debug("hyperweb application is installed and registered, nothing to do")
		} else {
			response, err := a.register(*name)
			if err != nil {
				return err
			}
			if response.Success {
				a.logger.Info("successfully registered cluster with gateway", zap.String("cluster name", *name))
			} else {
				return fmt.Errorf("failed to register cluster %s with gateway", *name)
			}
		}
	} else {
		a.logger.Info("hyperweb application is not installed - will attempt to install now")

		err = a.installHyperWeb(*name)
		if err != nil {
			return fmt.Errorf("failed to install hyperweb application: %v", err)
		}
	}

	return nil
}
