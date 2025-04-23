package hyperweb

import (
	"time"

	"epitome.hyperbolic.xyz/config"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func RunLoop(
	cfg config.Config,
	logger *zap.Logger,
	clientset kubernetes.Clientset,
	dynamicClient dynamic.DynamicClient,
	interval time.Duration,
) {
	ticker := time.NewTicker(interval)
	for {
		err := reconcile(clientset, dynamicClient, cfg.HYPERBOLIC_GATEWAY_URL, cfg.HYPERBOLIC_TOKEN)
		if err != nil {
			logrus.Fatalf("failed to reconcile: %v", err)
		}

		<-ticker.C
	}
}
