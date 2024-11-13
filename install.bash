#!/bin/bash

set -e

TOKEN=$TOKEN
EXTRA_PARAMS=""

if [[ "$DEV" == "true" ]]; then
  EXTRA_PARAMS="--set ref=dev"
fi

sudo snap install microk8s --classic --channel=1.31

# https://microk8s.io/docs/how-to-ceph
# https://canonical-microceph.readthedocs-hosted.com/en/reef-stable/tutorial/single-node/
sudo snap install microceph
sudo microceph cluster bootstrap


# TODO how much disk to create?
# sudo microceph disk add loop,4G,3

# may need to add:
# sudo modprobe rbd

# https://docs.ceph.com/en/reef/
# TODO we probably want pool size to be 1 for our use-case, or else storage
# will be duplicated without reason.
# sudo microceph.ceph config set global osd_pool_default_size 2
# sudo microceph.ceph config set mgr mgr_standby_modules false
# sudo microceph.ceph config set osd osd_crush_chooseleaf_type 0


# wait for this to be "HEALTH_OK"
sudo microceph.ceph status


echo "Starting microk8s..."
sudo microk8s start

echo "Enabling microk8s components..."
sudo microk8s enable rbac
sudo microk8s enable community
sudo microk8s enable argocd
sudo microk8s enable nvidia

sudo microk8s enable rook-ceph
sudo microk8s connect-external-ceph

sudo microk8s kubectl create namespace hyperdos
sudo microk8s kubectl create namespace hyperweb
sudo microk8s kubectl create namespace instance 
sudo microk8s kubectl create namespace ping

sudo microk8s helm repo add hyperdos https://hyperboliclabs.github.io/Hyper-dOS

echo "Starting hyperdos software..."
if [[ "$DEV" == "true" ]]; then
  echo "Install in dev mode..."
fi

namespace="argocd"

echo "Waiting for $namespace components to be ready..."

while true; do
  pods=$(sudo microk8s kubectl get pods -n "$namespace" --no-headers 2>&1)

  if [ "$pods" == "No resources found in $namespace namespace." ]; then
    echo "All $namespace components not ready yet."
    sleep 5
    continue
  fi

  # Get the number of non-ready pods
  not_ready_pods=$(echo $pods | awk '$3 != "Running" && $3 != "Completed" {print $1}' | wc -l)

  if [ "$not_ready_pods" -eq 0 ]; then
    echo "All $namespace components are ready!"
    break
  else
    echo "$not_ready_pods $namespace components are not ready yet. Checking again in 10 seconds..."
    sleep 10
  fi
done

echo "Installing hyperdos now..."
sleep 20 

sudo microk8s helm install hyperdos hyperdos/hyperdos --version 0.0.1-alpha.4 --set token=$TOKEN $EXTRA_PARAMS
