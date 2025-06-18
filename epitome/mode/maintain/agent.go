package maintain

import (
	"epitome.hyperbolic.xyz/config"
	argo "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type agent struct {
	cfg           config.Config
	logger        *zap.Logger
	clientset     kubernetes.Interface
	dynamicClient *dynamic.DynamicClient
	argoClient    argo.Interface
}
