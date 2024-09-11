# Hyperbolic Distributed Operating System

## Prerequisites

- You will need a Kubernetes cluster with argocd, and the NVIDIA operator installed.
- You will need to create the `hyperdos` namespace in your cluster

### install microk8s (Linux)

<https://microk8s.io/docs/getting-started>

``` shell
sudo snap install microk8s --classic --channel=1.30
sudo usermod -a -G microk8s $USER # allow non-sudo use of microk8s command
newgrp microk8s # reload shell

microk8s start # boot the cluster
microk8s enable rbac # improve security
sudo microk8s enable community # add the community repos
microk8s enable argocd # install argocd
microk8s enable nvidia # install the nvidia GPU operator

microk8s kubectl create namespace hyperdos # create the necessary namespace
```

### (optional) add more nodes to your cluster

<https://microk8s.io/docs/clustering>

1. (on the new node) `sudo snap install microk8s --classic --channel=1.30`
2. (on the original node) `microk8s add-node`
3. (on the new node) `microk8s join <output-from-original-node>`

## Install HyperdOS

1. login to <https://app.hyperbolic.xyz> and select 'settings'

2. copy your API Key and make sure to insert it in place of "<YOUR_API_KEY>" to run the installation command below:

``` shell
HYPERBOLIC_API_KEY=<YOUR_API_KEY> \
   && curl https://raw.githubusercontent.com/HyperbolicLabs/Hyper-dOS/main/install.yaml \
   | sed -e "s;{{stand-in}};${HYPERBOLIC_API_KEY};g" \
      | microk8s kubectl apply -f -
```

# Notes

- you only have to run this command on one node, and all your nodes will be added to the hyperweb

- if you already have nvidia drivers and container toolkit installed, use this command instead of `microk8s enable gpu`:

``` shell
microk8s enable nvidia --gpu-operator-driver host
```

  - you can override more NVIDIA GPU Operator settings by using the `--values` flag and referring to the values.yaml file here:
    - <https://github.com/NVIDIA/gpu-operator/blob/master/deployments/gpu-operator/values.yaml>

  - see further configuration options here:
    - <https://microk8s.io/docs/addon-gpu>

  - if the driver-validation container fails with this message:

``` shell
    Error: error validating driver installation:
     error creating symlink creator:
      failed to create NVIDIA device nodes:
       failed to create device node nvidiactl:
        failed to determine major:
         invalid device node

    Failed to create symlinks under /dev/char that point to all possible NVIDIA charact er devices. The existence of these symlinks is required to address the following bug: https://github.com/NVIDIA/gpu-operator/issues/430 bug impacts container runtimes configured with systemd cgroup management enabled.

    To disable the symlink creation, set the following envvar in ClusterPolicy:
```

Try creating the relevant envvar in the nvidia ClusterPolicy resource. There's a good chance the system will boot and operate normally. Edit the Nvidia ClusterPolicy validator.driver.env like so:

``` shell
        apiVersion: nvidia.com/v1
        kind: ClusterPolicy
...
        validator:
          driver:
            env:
            - name: DISABLE_DEV_CHAR_SYMLINK_CREATION
              value: "true"
```

- we do not officially support operating systems other than Linux. That being said, if you would like to join the Hyperbolic Supply Network from a Windows or MacOS device, you are welcome to give it a shot:
  - <https://microk8s.io/docs/install-alternatives>


- While most properly configured Kubernetes clusters should be able to run HyperdOS, for single-node clusters we officially support the microk8s distro only.


# Customized installation

If you would like to apply the installation manifest yourself rather than curling from github, you are welcome to download and edit the [install.yaml](install.yaml) file before applying it to your cluster


# (unstable preview) helm installation


## configure helm repo and dry-run
```shell
sudo microk8s helm repo add hyperdos https://hyperboliclabs.github.io/Hyper-dOS
sudo microk8s helm install --dry-run hyperdos hyperdos/hyperdos --version 0.0.1-alpha.4 --set ref="main" --set token="DRY_RUN_NO_TOKEN"
```


## install
``` shell
sudo microk8s helm install hyperdos hyperdos/hyperdos --version 0.0.1-alpha.4 token="<YOUR_API_KEY>"
```

## install (without rolling updates)
``` shell
# to disable automatic updates and pin to a specific git ref
sudo microk8s helm install hyperdos hyperdos/hyperdos --version 0.0.1-alpha.4 --set ref="0.0.1-alpha.4" token="<YOUR_API_KEY>"
```

## uninstall hyperdos
```shell
# to uninstall
sudo microk8s helm uninstall hyperdos
sudo microk8s kubectl delete app hyperweb -n argocd
```
