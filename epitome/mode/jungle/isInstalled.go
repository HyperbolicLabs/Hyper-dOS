package jungle

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func HyperwebIsInstalled(dynamicClient dynamic.DynamicClient) bool {
	_, err := dynamicClient.
		Resource(argoAppGVR).
		Namespace(argocdNamespace).
		Get(
			context.TODO(),
			hyperwebNamespace,
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
