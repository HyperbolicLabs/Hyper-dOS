package maintain

import (
	"fmt"
	"time"

	"epitome.hyperbolic.xyz/config"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func Run(
	cfg config.Config,
	logger *zap.Logger,
	clientset kubernetes.Interface,
	dynamicClient *dynamic.DynamicClient,
) error {

	a := &agent{
		cfg:           cfg,
		logger:        logger,
		clientset:     clientset,
		dynamicClient: dynamicClient,
	}

	interval := a.cfg.Maintain.ReconcileInterval
	a.logger.Info("running maintainance agent", zap.String("interval", interval.String()))

	ticker := time.NewTicker(interval)
	for {
		// patch cluster policy if we are on a buffalo baron
		// (as only the buffalo are expected to have the NVIDIA operator installed)
		if a.cfg.Role.Buffalo {
			err := a.patchClusterPolicy()
			if err != nil {
				logrus.Errorf("failed to patch cluster policy: %v", err)
				return err
			}
		}

		<-ticker.C // in maintain mode, we wait before running the first reconcile
		err := a.reconcile()
		if err != nil {
			return fmt.Errorf("failed to reconcile: %v", err)
		}
	}
}
