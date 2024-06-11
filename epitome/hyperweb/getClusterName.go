package hyperweb

import (
	"fmt"

	"epitome.hyperbolic.xyz/helper"
	"k8s.io/client-go/kubernetes"
)

func GetClusterName(clientSet kubernetes.Clientset) (*string, error) {

	cm, err := helper.GetConfigMap(clientSet, hyperwebNamespace, "cluster-name")
	if err != nil {
		return nil, err
	}

	if cm == nil {
		err = fmt.Errorf("cluster-name configmap is nil")
		return nil, err
	}

	if cm.Data == nil {
		err = fmt.Errorf("cluster-name configmap data is nil")
		return nil, err
	}

	// check if clusterName is set in configmap data
	if cm.Data["clusterName"] == "" {
		err = fmt.Errorf("cluster-name configmap data is empty")
		return nil, err
	}

	name := cm.Data["clusterName"]
	return &name, nil
}
