#!/bin/bash

set -e

TOKEN=$TOKEN

sudo snap install microk8s --classic --channel=1.31

echo "Starting microk8s..."
sudo microk8s start

echo "Enabling microk8s components..."
sudo microk8s enable rbac
sudo microk8s enable community
sudo microk8s enable argocd
sudo microk8s enable nvidia

sudo microk8s kubectl create namespace hyperdos
sudo microk8s helm repo add hyperdos https://hyperboliclabs.github.io/Hyper-dOS

echo "Starting hyperdos software..."
sudo microk8s helm install hyperdos hyperdos/hyperdos --version 0.0.1-alpha.4 --set token=$TOKEN
