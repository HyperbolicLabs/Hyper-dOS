package maintain

import (
	"fmt"
	"time"

	"epitome.hyperbolic.xyz/config"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

func Run(
	logger *zap.Logger,
	cfg *config.Config,
	clientset kubernetes.Interface,
) error {

	a := &agent{
		cfg:       *cfg,
		logger:    logger,
		clientset: clientset,
	}

	a.logger.Info("running maintainance agent")

	ticker := time.NewTicker(a.cfg.Maintain.ReconcileInterval)
	for {
		<-ticker.C // in maintain mode, we wait before running the first reconcile
		err := a.reconcile()
		if err != nil {
			return fmt.Errorf("failed to reconcile: %v", err)
		}
	}
}
