package jungle

import (
	"epitome.hyperbolic.xyz/config"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type agent struct {
	cfg           config.Config
	logger        *zap.Logger
	clientset     kubernetes.Interface
	dynamicClient *dynamic.DynamicClient
}
