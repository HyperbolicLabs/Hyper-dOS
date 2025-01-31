# Hyperbolic Distributed Operating System


## Install HyperdOS (single-node setup from scratch)

[ TODO curl -fsSL https://install.hyperbolic.xyz/hyperdos.sh | bash ]: #

```bash
curl -o install.bash https://raw.githubusercontent.com/HyperbolicLabs/Hyper-dOS/refs/heads/main/install.bash && bash install.bash
```

### Notes

- If you would like to run the install script yourself rather than curling from github, you are welcome to download and edit the [install.bash](https://github.com/HyperbolicLabs/Hyper-dOS/blob/main/install.bash) file before running it on your node.

- We do not officially support operating systems other than Linux. That being said, if you would like to join the Hyperbolic Supply Network from a Windows or MacOS device, you are welcome to give it a shot:

  - <https://microk8s.io/docs/install-alternatives>

- We officially support single-node microk8s+microceph clusters only, HOWEVER - a custom multi-node cluster should work smoothly if configured properly. See below for customized installation guidelines:



## (experimental) add more nodes to your cluster

Note: please reach out before doing this so we can support your cluster smoothly. As of right now, only single-node clusters will be shown to renters automatically.


### add microk8s node
<https://microk8s.io/docs/clustering>

1. (on the new node) `sudo snap install microk8s --classic --channel=1.32`
2. (on the original node) `microk8s add-node`
3. (on the new node) `microk8s join <output-from-original-node>`

### add microceph node

<https://microk8s.io/docs/clustering>

1. (on the new node) `sudo snap install microceph`
2. (on the new node) `sudo microceph init`



## Customized installation (existing kubernetes cluster)

Please get in touch if you are planning to install hyperdos on an existing multi-node cluster, we can help you get set up smoothly.

### Prerequisites
- ArgoCD installed: <https://argo-cd.readthedocs.io/en/stable/operator-manual/installation/>
- NVIDIA Operator installed: <https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/latest/getting-started.html>
- Namespaces `hyperdos` and `instance`
- You will need a `StorageClass` for rental instances to create PersistentVolumeClaims. We recommend `microceph`: <https://github.com/canonical/microceph>
- A ResourceQuota named `hyperstore` in the `instance` namespace. This will designate how much storage the network can use on your cluster.
- Please ensure at least 150GB of free disk space on each node before installing HyperdOS. Low disk space may lead to issues with your cluster, and failed rentals.


### configure helm repo and dry-run

```bash
sudo microk8s helm install --dry-run hyperdos hyperdos/hyperdos --version 0.0.1-beta.3 --set token="DRY_RUN_NO_TOKEN"
```

### install (without rolling updates)

```bash
# to disable automatic updates and pin to a specific git ref
sudo microk8s helm install hyperdos hyperdos/hyperdos --version 0.0.1-beta.3 --set token="<YOUR_API_KEY>"
```


## Uninstall

### remove hyperdos from the cluster

```bash
# to uninstall
sudo microk8s helm uninstall hyperdos
sudo microk8s kubectl delete app hyperweb -n argocd
```

### delete the cluster entirely

``` bash
# run these commands on each node in your cluster
sudo snap remove --purge microk8s
sudo snap remove --purge microceph
```


