# Hyperbolic Distributed Operating System

## Prerequisites

- You will need a Kubernetes cluster with argocd, and the NVIDIA operator installed.
- You will need to create the `hyperdos` namespace in your cluster

## Install HyperdOS

1. Login to <https://app.hyperbolic.xyz> and select 'settings'

2. Copy your API Key and make sure to insert it in place of "<YOUR_API_KEY>" to run the installation command below:

```bash
curl -o install.bash https://raw.githubusercontent.com/HyperbolicLabs/Hyper-dOS/refs/heads/dev/install.bash && chmod +x install.bash
TOKEN=<YOUR_API_KEY> ./install.bash
```

### (optional) add more nodes to your cluster

<https://microk8s.io/docs/clustering>

1. (on the new node) `sudo snap install microk8s --classic --channel=1.31`
2. (on the original node) `microk8s add-node`
3. (on the new node) `microk8s join <output-from-original-node>`

# Notes

- You only have to run this command on one node, and all your nodes will be added to the hyperweb

- We do not officially support operating systems other than Linux. That being said, if you would like to join the Hyperbolic Supply Network from a Windows or MacOS device, you are welcome to give it a shot:

  - <https://microk8s.io/docs/install-alternatives>

- While most properly configured Kubernetes clusters should be able to run HyperdOS, for single-node clusters we officially support the microk8s distro only.

# Customized installation

If you would like to apply the installation manifest yourself rather than curling from github, you are welcome to download and edit the [install.yaml](install.yaml) file before applying it to your cluster

## configure helm repo and dry-run

```bash
sudo microk8s helm install --dry-run hyperdos hyperdos/hyperdos --version 0.0.1-alpha.4 --set ref="main" --set token="DRY_RUN_NO_TOKEN"
```

## install (without rolling updates)

```bash
# to disable automatic updates and pin to a specific git ref
sudo microk8s helm install hyperdos hyperdos/hyperdos --version 0.0.1-alpha.4 --set ref="0.0.1-alpha.4" --set token="<YOUR_API_KEY>"
```

## uninstall hyperdos

```bash
# to uninstall
sudo microk8s helm uninstall hyperdos
sudo microk8s kubectl delete app hyperweb -n argocd
```
