package monkey

import (
	"time"

	"epitome.hyperbolic.xyz/cluster"
	"epitome.hyperbolic.xyz/config"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type agent struct {
	cfg           config.Config
	logger        *zap.Logger
	clientset     kubernetes.Clientset
	dynamicClient dynamic.DynamicClient
}

func Run(
	cfg config.Config,
	logger *zap.Logger,
) error {
	clientset, dynamicClient := cluster.MustConnect(cfg.KUBECONFIG)

	a := &agent{
		cfg:           cfg,
		logger:        logger,
		clientset:     clientset,
		dynamicClient: dynamicClient,
	}

	ticker := time.NewTicker(1 * time.Minute)
	for {
		err := a.reconcile()
		if err != nil {
			logger.Error("reconcile failed", zap.Error(err))
		}

		<-ticker.C
	}
}
