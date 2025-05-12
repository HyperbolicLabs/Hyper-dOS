package cluster

import (
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func MustConnect(
	logger *zap.Logger,
	kubeconfig string,
	failoverToDefaultPath bool,
) (
	kubernetes.Interface,
	*dynamic.DynamicClient) {
	if kubeconfig != "" {
		// create the config object from kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logger.Fatal("failed to create clientconfig", zap.Error(err))
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			logger.Fatal("failed to create clientset", zap.Error(err))
		}

		dynamicClient, err := dynamic.NewForConfig(config)
		if err != nil {
			logger.Fatal("failed to create dynamic clientset", zap.Error(err))
		}

		logger.Info("connected to cluster", zap.String("kubeconfig", kubeconfig))

		return clientset, dynamicClient
	}

	if failoverToDefaultPath {
		logger.Info("using default kubeconfig", zap.String("kubeconfig", clientcmd.RecommendedHomeFile))
		config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			logger.Info("couldn't create config from kubeconfig file, will try in-cluster", zap.Error(err))
		} else {
			logger.Info("connected to cluster", zap.String("kubeconfig", clientcmd.RecommendedHomeFile))
			return kubernetes.NewForConfigOrDie(config), dynamic.NewForConfigOrDie(config)
		}
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Fatal("failed to create in-cluster config", zap.Error(err))
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatal("failed to create clientset", zap.Error(err))
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		logger.Fatal("failed to create dynamic clientset", zap.Error(err))
	}

	logger.Info("connected to cluster", zap.String("kubeconfig", "in-cluster"))
	return clientset, dynamicClient
}
