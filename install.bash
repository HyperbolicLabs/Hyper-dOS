#!/bin/bash

set -e

TOKEN=$TOKEN
EXTRA_PARAMS=""

if [[ "$DEV" == "true" ]]; then
  EXTRA_PARAMS="--set ref=dev"
fi

sudo snap install microk8s --classic --channel=1.31

echo "Starting microk8s..."
sudo microk8s start

echo "Enabling microk8s components..."
sudo microk8s enable rbac
sudo microk8s enable community
sudo microk8s enable argocd
sudo microk8s enable nvidia

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

sleep 15

echo "Installing hyperdos now..."
sudo microk8s helm install hyperdos hyperdos/hyperdos --version 0.0.1-alpha.4 --set token=$TOKEN $EXTRA_PARAMS
