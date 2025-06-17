package maintain

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func (a *agent) checkCephClusterHealth(
	ctx context.Context,
	namespacedName types.NamespacedName,
) error {

	// use the dynamic client to get the health of the ceph cluster
	// we should get the 'rook-ceph-external' cephcluster
	// from the 'rook-ceph-external' namespace by default

	unstructuredCephCluster, err := a.dynamicClient.Resource(
		schema.GroupVersionResource{
			Group:    "ceph.rook.io",
			Version:  "v1",
			Resource: "cephclusters",
		},
	).
		Namespace(namespacedName.Namespace).
		Get(ctx, namespacedName.Name, metav1.GetOptions{})

	if err != nil {
		return err
	}

	if unstructuredCephCluster == nil {
		return fmt.Errorf("cephcluster is nil")
	}

	// parse status.ceph.health from unstructured
	var limitedCephCluster limitedCephCluster
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(
		unstructuredCephCluster.Object,
		&limitedCephCluster)

	if limitedCephCluster.Status.Ceph.Health != "HEALTH_OK" {
		// check if HEALTH_WARN
		if limitedCephCluster.Status.Ceph.Health == "HEALTH_WARN" {
			// check if the reason is because too little replication
			// this is the default for single-node hyperdos clusters,
			// so we can ignore it (for now, down the line there can be improvements)
			if strings.Contains(limitedCephCluster.Status.Ceph.Details.POOL_NO_REDUNDANCY.Message, "no replicas configured") {
				a.logger.Debug("skipping ceph replication warning",
					zap.Any("details", limitedCephCluster.Status.Ceph.Details))
			} else {
				return fmt.Errorf("ceph cluster: %v", limitedCephCluster.Status.Ceph.Health)
			}
		} else {
			return fmt.Errorf("ceph cluster: %v", limitedCephCluster.Status.Ceph.Health)
		}
	}

	return nil
}

const (
	cephClusterGVR = "ceph.rook.io/v1"
)

type limitedCephCluster struct {
	Status struct {
		Ceph struct {
			Health  string `json:"health"`
			Details struct {
				POOL_NO_REDUNDANCY struct {
					Message  string `json:"message"`
					Severity string `json:"severity"`
				} `json:"POOL_NO_REDUNDANCY"`
			} `json:"details"`
		} `json:"ceph"`
	} `json:"status"`
}

const ex = `
apiVersion: ceph.rook.io/v1
kind: CephCluster
metadata:
  annotations:
    meta.helm.sh/release-name: rook-ceph-external
    meta.helm.sh/release-namespace: rook-ceph-external
  creationTimestamp: "2024-12-02T19:19:23Z"
  finalizers:
  - cephcluster.ceph.rook.io
  generation: 2
  labels:
    app.kubernetes.io/managed-by: Helm
  name: rook-ceph-external
  namespace: rook-ceph-external
  resourceVersion: "40111790"
  uid: 8993261d-f1e5-4aa7-90e6-d3fe672647f6
spec:
  cephVersion:
    image: quay.io/ceph/ceph:v18.2.4
  cleanupPolicy:
    sanitizeDisks:
      dataSource: zero
      iteration: 1
      method: quick
  crashCollector:
    disable: true
  dashboard:
    enabled: true
    ssl: true
  dataDirHostPath: /var/lib/rook
  disruptionManagement:
    managePodBudgets: true
    osdMaintenanceTimeout: 30
  external:
    enable: true
  healthCheck:
    daemonHealth:
      mon:
        interval: 45s
      osd:
        interval: 1m0s
      status:
        interval: 1m0s
    livenessProbe:
      mgr: {}
      mon: {}
      osd: {}
  logCollector:
    enabled: true
    maxLogSize: 500M
    periodicity: daily
  mgr:
    count: 2
  mon:
    count: 3
  monitoring: {}
  network:
    connections:
      compression: {}
      encryption: {}
    multiClusterService: {}
  priorityClassNames:
    mgr: system-cluster-critical
    mon: system-node-critical
    osd: system-node-critical
  resources:
    cleanup:
      limits:
        memory: 1Gi
      requests:
        cpu: 500m
        memory: 100Mi
    crashcollector:
      limits:
        memory: 60Mi
      requests:
        cpu: 100m
        memory: 60Mi
    exporter:
      limits:
        memory: 128Mi
      requests:
        cpu: 50m
        memory: 50Mi
    logcollector:
      limits:
        memory: 1Gi
      requests:
        cpu: 100m
        memory: 100Mi
    mgr:
      limits:
        memory: 1Gi
      requests:
        cpu: 500m
        memory: 512Mi
    mgr-sidecar:
      limits:
        memory: 100Mi
      requests:
        cpu: 100m
        memory: 40Mi
    mon:
      limits:
        memory: 2Gi
      requests:
        cpu: "1"
        memory: 1Gi
    osd:
      limits:
        memory: 4Gi
      requests:
        cpu: "1"
        memory: 4Gi
    prepareosd:
      requests:
        cpu: 500m
        memory: 50Mi
  security:
    keyRotation:
      enabled: false
    kms: {}
  skipUpgradeChecks: true
  storage:
    useAllDevices: true
    useAllNodes: true
  waitTimeoutForHealthyOSDInMinutes: 10
status:
  ceph:
    capacity:
      bytesAvailable: 268070887424
      bytesTotal: 268435456000
      bytesUsed: 364568576
      lastUpdated: "2025-06-17T22:22:52Z"
    details:
      POOL_NO_REDUNDANCY:
        message: 2 pool(s) have no replicas configured
        severity: HEALTH_WARN
    fsid: df4f9b3a-49dc-4921-b60e-db490ba440e9
    health: HEALTH_WARN
    lastChecked: "2025-06-17T22:22:52Z"
    versions:
      mds:
        ceph version 19.2.0 (16063ff2022298c9300e49a547a16ffda59baf13) squid (stable): 1
      mgr:
        ceph version 19.2.0 (16063ff2022298c9300e49a547a16ffda59baf13) squid (stable): 1
      mon:
        ceph version 19.2.0 (16063ff2022298c9300e49a547a16ffda59baf13) squid (stable): 1
      osd:
        ceph version 19.2.0 (16063ff2022298c9300e49a547a16ffda59baf13) squid (stable): 1
      overall:
        ceph version 19.2.0 (16063ff2022298c9300e49a547a16ffda59baf13) squid (stable): 4
  conditions:
  - lastHeartbeatTime: "2025-06-17T15:29:02Z"
    lastTransitionTime: "2025-06-17T15:29:02Z"
    message: Attempting to connect to an external Ceph cluster
    reason: ClusterConnecting
    status: "True"
    type: Connecting
  - lastHeartbeatTime: "2025-06-17T22:22:53Z"
    lastTransitionTime: "2024-12-02T19:19:32Z"
    message: Cluster connected successfully
    reason: ClusterConnected
    status: "True"
    type: Connected
  message: Cluster connected successfully
  phase: Connected
  state: Connected
  version:
    image: quay.io/ceph/ceph:v18.2.4
    version: 18.2.4-0
`
