package hyperweb

import (
	"fmt"
	"time"

	"epitome.hyperbolic.xyz/config"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func RunLoop(
	cfg config.Config,
	logger *zap.Logger,
	clientset kubernetes.Clientset,
	dynamicClient dynamic.DynamicClient,
) error {
	interval := cfg.Default.ReconcileInterval
	ticker := time.NewTicker(interval)
	for {
		err := reconcile(clientset, dynamicClient, cfg.Default.HYPERBOLIC_GATEWAY_URL, cfg.Default.HYPERBOLIC_TOKEN)
		if err != nil {
			return fmt.Errorf("default epitome failed to reconcile: %v", err)
		}

		<-ticker.C
	}
}
