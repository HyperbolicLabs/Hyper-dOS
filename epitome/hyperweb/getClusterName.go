package hyperweb

import (
	"fmt"

	"epitome.hyperbolic.xyz/helper"
	"k8s.io/client-go/kubernetes"
)

func GetClusterName(clientSet kubernetes.Interface) (*string, error) {

	cm, err := helper.GetConfigMap(clientSet, hyperdosNamespace, "cluster-name")
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

	name := cm.Data[clusterNameDataField]
	// check if clusterName is set in configmap data
	if name == "" {
		err = fmt.Errorf("cluster-name configmap data is empty for field %v", clusterNameDataField)
		return nil, err
	}

	return &name, nil
}
