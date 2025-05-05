package jungle

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *agent) HyperwebIsInstalled() bool {
	_, err := a.dynamicClient.
		Resource(argoAppGVR).
		Namespace(argocdNamespace).
		Get(
			context.TODO(),
			a.cfg.HyperwebNamespace,
			metav1.GetOptions{})

	if err != nil {
		// logrus.Errorf("list: %v, err: %v", appset, err)
		// panic(err)
		return false
	}

	// print the unstructured object in json format
	// logrus.Infof("unstructured object: %v", us)

	return true
}

var argoAppGVR = schema.GroupVersionResource{
	Group:    "argoproj.io",
	Version:  "v1alpha1",
	Resource: "applications",
}
