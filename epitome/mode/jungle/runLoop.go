package jungle

import (
	"fmt"
	"time"

	"epitome.hyperbolic.xyz/config"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func Run(
	cfg config.Config,
	logger *zap.Logger,
	clientset kubernetes.Interface,
	dynamicClient dynamic.DynamicClient,
) error {
	a := &agent{
		cfg:       cfg,
		logger:    logger,
		clientset: clientset,
	}
	interval := cfg.Default.ReconcileInterval
	ticker := time.NewTicker(interval)
	for {
		err := a.reconcile()
		if err != nil {
			return fmt.Errorf("default epitome failed to reconcile: %v", err)
		}

		<-ticker.C
	}
}
