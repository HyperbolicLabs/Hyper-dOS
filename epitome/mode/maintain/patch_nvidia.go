package maintain

import (
	"context"
	"time"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func (a *agent) patchClusterPolicy() error {
	gvr := schema.GroupVersionResource{
		Group:    "nvidia.com",
		Version:  "v1",
		Resource: "clusterpolicies",
	}

	parentCtx := context.Background()

	getCtx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	clusterPolicy, err := a.dynamicClient.Resource(gvr).Get(getCtx, "cluster-policy", metav1.GetOptions{})
	if err != nil {
		a.logger.Error("failed to get ClusterPolicy", zap.Error(err))
		return err
	}

	_, found, err := unstructured.NestedMap(clusterPolicy.Object, "spec", "validator", "driver")
	if err != nil {
		a.logger.Error("error checking .spec.validator.driver", zap.Error(err))
		return err
	}

	if !found {
		patchCtx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
		defer cancel()

		patch := []byte(`{"spec":{"validator":{"driver":{"env":[{"name":"DISABLE_DEV_CHAR_SYMLINK_CREATION","value":"true"}]}}}}`)
		_, err = a.dynamicClient.Resource(gvr).Patch(patchCtx, "cluster-policy", types.MergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			a.logger.Error("failed to patch ClusterPolicy", zap.Error(err))
			return err
		}
		a.logger.Info("successfully patched ClusterPolicy with .spec.validator.driver.env")
	} else {
		a.logger.Debug(".spec.validator.driver field already exists in ClusterPolicy")
	}

	return nil
}
