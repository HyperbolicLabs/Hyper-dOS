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

sudo microk8s helm install hyperdos hyperdos/hyperdos --version 0.0.1-alpha.4 --set token=$TOKEN $EXTRA_PARAMS
