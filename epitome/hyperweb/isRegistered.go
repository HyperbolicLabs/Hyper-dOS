package hyperweb

import (
	"epitome.hyperbolic.xyz/helper"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

func isRegistered(
	clientset kubernetes.Clientset,
) bool {
	// check if ping configmap exists in ping namespace
	cm, err := helper.GetConfigMap(clientset, pingNamespace, "ping")

	// if doesn't exist
	if errors.IsNotFound(err) {
		return false
	}

	if err != nil {
		logrus.Fatalf("unexpected error trying to get ping configmap: %v", err)
	}

	if cm == nil {
		logrus.Fatalf("ping configmap is nil")
	}

	if cm.Data == nil {
		logrus.Fatalf("ping configmap data is nil")
	}

	value, ok := cm.Data["ping"]

	if !ok {
		logrus.Fatalf("ping configmap data is empty for field ping")
	}

	if value != "pong" {
		logrus.Fatalf("ping configmap data is not pong")
	}

	return true
}
