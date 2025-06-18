package maintain

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *agent) checkDaemonsetsHealthy(
	ctx context.Context,
	namespace string,
) error {
	daemonsets, err := a.clientset.AppsV1().
		DaemonSets(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	// if no daemonsets, no problem
	for _, ds := range daemonsets.Items {
		if ds.Status.NumberUnavailable > 0 {
			a.logger.Warn(
				"daemonset unhealthy",
				zap.String("daemonset", ds.Name),
				zap.String("namespace", namespace),
				zap.Int32("unavailable", ds.Status.NumberUnavailable),
			)
			return fmt.Errorf(
				"daemonset %s in namespace %s is not ready: %v unavailable",
				ds.Name,
				namespace,
				ds.Status.NumberUnavailable,
			)
		}
	}

	return nil
}
