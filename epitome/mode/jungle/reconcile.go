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
			"operator-oauth secret does not exist",
			zap.String("namespace", a.cfg.HyperwebNamespace))

		response, err := a.handshake(http.DefaultClient)
		if err != nil {
			a.logger.Error("failed to handshake with gateway", zap.Error(err))
			return err
		}

		a.mustCreateTailscaleOperatorOAuthSecret(
			response.ClientID,
			response.ClientSecret,
		)

		err = a.installClusterNameConfigMap(response.ClusterName)
		if err != nil {
			a.logger.Error("failed to save cluster name in configmap: %v", zap.Error(err))
			return err
		}
	}

	name, err := a.getClusterName()
	if err != nil {
		a.logger.Error("failed to get cluster name: %v", zap.Error(err))
		return err
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
				a.logger.Info("registered cluster with gateway", zap.String("cluster name", *name))
			} else {
				return fmt.Errorf("failed to register cluster %s with gateway", *name)
			}
		}
	} else {
		a.logger.Info("hyperweb application is not installed - installing now")

		err = a.installHyperWeb(*name)
		if err != nil {
			a.logger.Error("failed to install hyperweb application", zap.Error(err))
			return err
		}
	}
	return nil
}
