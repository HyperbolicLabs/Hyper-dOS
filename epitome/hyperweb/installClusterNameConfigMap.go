package hyperweb

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func installClusterNameConfigMap(
	clientset kubernetes.Clientset,
	clusterName string,
) error {

	configmap := &v1.ConfigMap{}
	configmap.Name = "cluster-name"
	configmap.Namespace = hyperdosNamespace
	configmap.Data = map[string]string{
		clusterNameDataField: clusterName,
	}

	cm, err := clientset.CoreV1().ConfigMaps(configmap.Namespace).Create(context.TODO(), configmap, metav1.CreateOptions{})

	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		} else {
			logrus.Debugf("configmap %s already exists", configmap.Name)
			// check if clusterName is set in configmap data

			value, ok := cm.Data[clusterNameDataField]
			if !ok {
				return fmt.Errorf("cluster-name configmap data is empty for field %v", clusterNameDataField)
			} else if value != clusterName {
				logrus.Warnf("configmap %s already exists with different %s value. existing value: %s attempted value: %s. Will not update.",
					configmap.Name,
					clusterNameDataField,
					value,
					clusterName)
			}
		}
	}

	return nil
}
