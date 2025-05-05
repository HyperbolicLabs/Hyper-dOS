package maintain

import (
	"context"
	"testing"
	"time"

	"epitome.hyperbolic.xyz/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestRestartCalicoIfExists(t *testing.T) {
	// Setup test environment
	ctx := context.Background()
	clientset := fake.NewSimpleClientset()
	testStartTime := time.Now().Add(-1 * time.Second)

	// Create agent with fake client
	a := &agent{
		clientset: clientset,
		logger:    zap.NewNop(),
	}

	// reconcile should not return an error if the daemonset does not exist
	err := a.restartCalicoIfExists(ctx)
	require.NoError(t, err)

	// Create calico daemonset
	_, err = clientset.AppsV1().DaemonSets(config.CalicoNamespace).Create(ctx, &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "calico-node",
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	// Create some dummy pods
	_, err = clientset.CoreV1().Pods(config.CalicoNamespace).Create(ctx, &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "calico-node-abc",
			Labels: map[string]string{
				"k8s-app": "calico-node",
			},
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	// Execute the restart
	err = a.restartCalicoIfExists(ctx)
	require.NoError(t, err)

	// check that the restartedat annotation was added in the proper format
	daemonset, err := clientset.AppsV1().DaemonSets(config.CalicoNamespace).Get(ctx, "calico-node", metav1.GetOptions{})
	require.NoError(t, err)
	restartedAt := daemonset.Annotations["kubectl.kubernetes.io/restartedAt"]
	require.NotEmpty(t, restartedAt)

	// restartedAt should be after test start time
	parsedRestartedAt, err := time.Parse(time.RFC3339, restartedAt)
	require.NoError(t, err)
	require.True(t, parsedRestartedAt.After(testStartTime))
}
