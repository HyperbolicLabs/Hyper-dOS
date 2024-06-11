package hyperweb

import (
	"bytes"
	"context"
	"html/template"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const yamlConfigMap = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-name
  namespace: {{.Namespace}}
data:
  {{.ClusterNameDataField}}: {{.ClusterName}}
`

var cmGVR = schema.GroupVersionResource{
	Group:    "",
	Version:  "v1",
	Resource: "configmaps",
}

func InstallCM(dynamicClient dynamic.DynamicClient, clusterName string) error {
	templatedcm := template.Must(template.New("index").Parse(yamlConfigMap))
	buff := new(bytes.Buffer)

	err := templatedcm.Execute(
		buff,
		map[string]interface{}{
			"Namespace":            hyperdosNamespace,
			"ClusterNameDataField": clusterNameDataField,
			"ClusterName":          clusterName,
		},
	)

	if err != nil {
		logrus.Fatal(err)
	}

	obj := &unstructured.Unstructured{}
	yaml.Unmarshal(buff.Bytes(), obj)

	us, err := dynamicClient.
		Resource(cmGVR).
		Namespace(hyperdosNamespace).
		Apply(
			context.TODO(),
			clusterNameDataField,
			obj,
			metav1.ApplyOptions{
				FieldManager: "epitome",
			})

	if err != nil {
		logrus.Errorf("could not apply yaml configmap: %v", err)
		panic(err)
	}

	logrus.Infof("applied unstructured: %+v", us)

	return nil
}
