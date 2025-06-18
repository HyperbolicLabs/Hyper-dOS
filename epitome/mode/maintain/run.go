package maintain

import (
	"fmt"
	"time"

	"epitome.hyperbolic.xyz/config"
	argo "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func Run(
	cfg config.Config,
	logger *zap.Logger,
	clientset kubernetes.Interface,
	dynamicClient *dynamic.DynamicClient,
	argoClient argo.Interface,
) error {

	a := &agent{
		cfg:           cfg,
		logger:        logger,
		clientset:     clientset,
		dynamicClient: dynamicClient,
		argoClient:    argoClient,
	}

	interval := a.cfg.Maintain.ReconcileInterval
	a.logger.Info("running maintenance agent", zap.String("interval", interval.String()))

	// hack: fix common canonical/calico bug
	go a.occasionallyRestartCalico()

	ticker := time.NewTicker(interval)
	for {
		err := a.reconcile()
		if err != nil {
			// break the loop
			return fmt.Errorf("failed to reconcile: %v", err)
		}

		<-ticker.C
	}
}
