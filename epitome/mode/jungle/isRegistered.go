package jungle

import (
	"epitome.hyperbolic.xyz/helper"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (a *agent) isRegistered() bool {
	// check if ping configmap exists in ping namespace
	cm, err := helper.GetConfigMap(a.clientset, pingNamespace, "ping")

	// if doesn't exist
	if errors.IsNotFound(err) {
		logrus.Warnf("king cluster ping configmap not found: %v", err)
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
