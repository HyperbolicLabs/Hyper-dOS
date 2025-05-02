package maintain

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (a *agent) restartCalicoIfExists(ctx context.Context) error {
	// restart the calico daemonset
	calicoDaemonSetName := "calico-node"
	namespace := "calico-system"

	// "kubectl rollout restart daemonset calico-node"
	// since 'rollout restart' is an unofficial procedure in the api,
	// this happens under the hood by patching the daemonset
	// with the annotation:
	// "kubectl.kubernetes.io/restartedAt":"<timestamp>"
	// https://stackoverflow.com/questions/61335318/how-to-restart-a-deployment-in-kubernetes-using-go-client
	patch := fmt.Appendf(nil,
		`{
				"metadata": {
					"annotations": {
						"kubectl.kubernetes.io/restartedAt": "%s"
					}
				}
			}`,
		time.Now().Format(time.RFC3339),
	)

	_, err := a.clientset.AppsV1().
		DaemonSets(namespace).
		Patch(
			ctx,
			calicoDaemonSetName,
			types.StrategicMergePatchType,
			patch,
			metav1.PatchOptions{})
	if err != nil {
		// if err is not found, return nil
		if err.Error() == "daemonsets.apps \"calico-node\" not found" {
			a.logger.Debug("calico daemonset not found, this baron likely doesn't have calico installed. Will move on.")
			return nil
		}
		return fmt.Errorf("failed to patch calico node daemonset: %v", err)
	}

	a.logger.Info("calico daemonset restarted")
	return nil
}
