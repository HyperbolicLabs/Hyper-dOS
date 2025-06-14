package maintain

import (
	"context"
	"fmt"
	"time"

	"epitome.hyperbolic.xyz/config"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// occaisionallyRestartCalico restarts the calico daemonset every interval period if it exists.
// This is a hack to work around a bug in calico where it will sometimes crash all containers in the cluster
// when a container is started or stopped. This seems to happen when microk8s is installed
// without also running "snap refresh --hold microk8s" to block autoupdates
// calico or canonical should fix whatever is causing them to bork clusters on upgrades.
// until then, the hack remains.
func (a *agent) occasionallyRestartCalico() {
	ticker := time.NewTicker(a.cfg.Maintain.CalicoRestartInterval)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := a.restartCalicoIfExists(ctx)
		if err != nil {
			a.logger.Warn("failed to restart calico", zap.Error(err))
		}

		<-ticker.C
	}
}

func (a *agent) restartCalicoIfExists(ctx context.Context) error {
	// restart the calico daemonset
	calicoDaemonSetName := "calico-node"
	namespace := config.CalicoNamespace

	// "kubectl rollout restart daemonset calico-node"
	// since 'rollout restart' is an unofficial procedure in the api,
	// this happens under the hood by patching the daemonset
	// with the annotation:
	// "kubectl.kubernetes.io/restartedAt":"<timestamp>"
	// https://stackoverflow.com/questions/61335318/how-to-restart-a-deployment-in-kubernetes-using-go-client
	patch := fmt.Appendf(nil,
		`{
			"spec": {
				"template": {
					"metadata": {
						"annotations": {
							"kubectl.kubernetes.io/restartedAt": "%s"
						}
					}
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
