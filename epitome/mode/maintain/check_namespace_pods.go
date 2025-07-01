package maintain

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *agent) checkNamespaceForUnhealthyPods(
	ctx context.Context,
	namespace string,
) error {
	pods, err := a.clientset.CoreV1().Pods(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to check namespace %s for unhealthy pods. err: %v",
			namespace,
			err)
	}
	var troubledPodNames []string
	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodRunning &&
			pod.Status.Phase != corev1.PodSucceeded {
			troubledPodNames = append(troubledPodNames, pod.Name)
		}
	}

	if len(troubledPodNames) > 0 {
		return fmt.Errorf("troubled pods in namespace %s: %v",
			namespace,
			troubledPodNames)
	}

	return nil
}
