package jungle

import (
	"bytes"
	"context"
	"html/template"
	"os"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	yaml "k8s.io/apimachinery/pkg/util/yaml"
)

func InstallHyperWeb(
	dynamicClient dynamic.DynamicClient,
	clusterName string,
) error {
	yamlFile, err := os.ReadFile("application.yaml")
	if err != nil {
		logrus.Fatalf("failed to read application.yaml: %v", err)
	}

	templatedapp := template.Must(template.New("index").Parse(string(yamlFile)))
	buff := new(bytes.Buffer)
	err = templatedapp.Execute(
		buff,
		map[string]interface{}{
			"ClusterName": clusterName})

	if err != nil {
		logrus.Fatal("failed to generate application template", err)
	}

	obj := &unstructured.Unstructured{}
	yaml.NewYAMLOrJSONDecoder(buff, buff.Len()).Decode(obj)
	logrus.Debugf("unstructured: %+v", obj)

	_, err = dynamicClient.
		Resource(argoGVR).
		Namespace("argocd").
		Apply(
			context.TODO(),
			"hyperweb",
			obj,
			metav1.ApplyOptions{
				FieldManager: "epitome",
			})

	if err != nil {
		logrus.Errorf("could not apply yaml application: %v", err)
		panic(err)
	}

	return nil
}

var argoGVR = schema.GroupVersionResource{
	Group:    "argoproj.io",
	Version:  "v1alpha1",
	Resource: "applications",
}
