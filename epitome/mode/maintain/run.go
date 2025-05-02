package maintain

import (
	"fmt"
	"time"

	"epitome.hyperbolic.xyz/config"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

func Run(
	cfg config.Config,
	logger *zap.Logger,
	clientset kubernetes.Interface,
) error {

	a := &agent{
		cfg:       cfg,
		logger:    logger,
		clientset: clientset,
	}

	interval := a.cfg.Maintain.ReconcileInterval
	a.logger.Info("running maintainance agent", zap.String("interval", interval.String()))

	ticker := time.NewTicker(interval)
	for {
		<-ticker.C // in maintain mode, we wait before running the first reconcile
		err := a.reconcile()
		if err != nil {
			return fmt.Errorf("failed to reconcile: %v", err)
		}
	}
}
