package monkey

import (
	"context"
	"encoding/json"
	"fmt"

	"epitome.hyperbolic.xyz/config"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (a *agent) reconcile() error {
	nodeName := a.cfg.Monkey.KUBERNETES_NODE_NAME

	// Get current node object
	node, err := a.clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get node: %w", err)
	}

	// get the cpu name and clock
	cpuLabels, err := a.getThisCPU()
	if err != nil {
		return fmt.Errorf("failed to get CPU name: %w", err)
	}

	newLabels := make(map[string]string)
	newLabels[config.CPUNameLabelKey] = cpuLabels.name

	// check if label(s) already exist on the node and is correct
	if labelsAreGood(node.Labels, newLabels) {
		a.logger.Debug("node already labeled correctly",
			zap.String("node", nodeName),
			zap.String(config.CPUNameLabelKey, node.Labels[config.CPUNameLabelKey]))
		return nil
	}

	// Update node labels with a patch
	patchBytes, err := json.Marshal(map[string]any{
		"metadata": map[string]any{
			"labels": newLabels,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create patch json: %w", err)
	}
	_, err = a.clientset.CoreV1().Nodes().Patch(
		context.TODO(),
		nodeName,
		types.MergePatchType,
		patchBytes,
		metav1.PatchOptions{})

	if err != nil {
		return fmt.Errorf("failed to update node labels: %w", err)
	}

	a.logger.Info("successfully labeled node",
		zap.String("node", nodeName),
		zap.String(config.CPUNameLabelKey, cpuLabels.name),
	)

	return nil
}

func labelsAreGood(existingLabels map[string]string, newLabels map[string]string) bool {
	for k, v := range newLabels {
		if v2, ok := existingLabels[k]; !ok || v != v2 {
			return false
		}
	}

	return true
}
