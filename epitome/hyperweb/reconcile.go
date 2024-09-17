package hyperweb

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func reconcile(
	clientset kubernetes.Clientset,
	dynamicClient dynamic.DynamicClient,
	gatewayUrl string,
	token string,
) error {
	if !secretExists(clientset, hyperwebNamespace, "operator-oauth") {
		// if it does not, query the gateway for oauth credentials using our token
		logrus.Infof("operator-oauth secret does not exist in namespace: %v", hyperwebNamespace)

		response, err := handshake(gatewayUrl, token)
		if err != nil {
			logrus.Errorf("failed to handshake with gateway: %v", err)
			return err
		}

		mustCreateOperatorOAuthSecret(
			clientset,
			hyperwebNamespace,
			"operator-oauth",
			response.ClientID,
			response.ClientSecret,
		)

		err = installClusterNameConfigMap(clientset, response.ClusterName)
		if err != nil {
			logrus.Errorf("failed to save cluster name in configmap: %v", err)
			return err
		}
	}

	name, err := GetClusterName(clientset)
	if err != nil {
		logrus.Errorf("failed to get cluster name: %v", err)
		return err
	}

	if IsInstalled(dynamicClient) {
		if isRegistered(clientset) {
			logrus.Infof("hyperweb application is installed and registered, nothing to do")
		} else {
			response, err := register(
				gatewayUrl,
				token,
				*name,
			)
			if err != nil {
				return err
			}
			if response.Success {
				logrus.Infof("registered cluster %s with gateway", *name)
			} else {
				return fmt.Errorf("failed to register cluster %s with gateway", *name)
			}
		}
	} else {
		logrus.Infof("hyperweb application is not installed - installing now")

		err = InstallHyperWeb(dynamicClient, *name)
		if err != nil {
			logrus.Errorf("failed to install hyperweb application: %v", err)
			return err
		}
	}

	parentCtx := context.Background()

	// Get the ClusterPolicy resource
	gvr := schema.GroupVersionResource{
		Group:    "nvidia.com",
		Version:  "v1",
		Resource: "clusterpolicies",
	}

	getCtx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	clusterPolicy, err := dynamicClient.Resource(gvr).Get(getCtx, "cluster-policy", metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("failed to get ClusterPolicy: %v", err)
		return err
	}

	// Check if .spec.validator.driver field exists
	_, found, err := unstructured.NestedMap(clusterPolicy.Object, "spec", "validator", "driver")
	if err != nil {
		logrus.Errorf("error checking .spec.validator.driver: %v", err)
		return err
	}

	if !found {
		patchCtx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
		defer cancel()

		// Patch the ClusterPolicy with the required field
		patch := []byte(`{"spec":{"validator":{"driver":{"env":[{"name":"DISABLE_DEV_CHAR_SYMLINK_CREATION","value":"true"}]}}}}`)
		_, err = dynamicClient.Resource(gvr).Patch(patchCtx, "cluster-policy", types.MergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			logrus.Errorf("failed to patch ClusterPolicy: %v", err)
			return err
		}
		logrus.Infof("successfully patched ClusterPolicy with .spec.validator.driver.env")
	} else {
		logrus.Infof(".spec.validator.driver field already exists in ClusterPolicy")
	}

	return nil
}
