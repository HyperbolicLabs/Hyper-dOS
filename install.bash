#!/bin/bash

set -e

# configuration
API_URL="https://api.hyperbolic.xyz"
HYPERDOS_VERSION=0.0.1-beta.5
MICROK8S_VERSION=1.32
TOKEN=$TOKEN
EXTRA_PARAMS=""
PATH=$PATH:/var/lib/snapd/snap/bin

if [[ "$DEV" == "true" ]]; then
  set -x
  EXTRA_PARAMS="--set hyperdos.ref=dev"
  API_URL="https://api.dev-hyperbolic.xyz"
fi

###
## some helper functions
###

check_for_linux() {
  ostype=$(uname -s)
  if [ "$ostype" != "Linux" ]; then
    echo "This script can only be run on Linux"
    cancel
  fi
}

check_for_snap() {
  if ! command -v snap >/dev/null 2>&1; then
    echo "snap is not installed"
    exit 1
  fi
}

check_for_token() {
  if [ -z "$TOKEN" ]; then
    echo "the TOKEN environment variable is not set"
    echo "your Hyperbolic API Key can be found at https://app.hyperbolic.xyz/settings"
    read -r -p "Please enter your API Key: " TOKEN
  fi
}

validate_token() {
  echo "querying Hyperbolic supply API to validate api key..."

  # curl -s -I -X \
  #   GET $API_URL/v1/marketplace/instances/supplied \
  #   -H "Authorization: Bearer $TOKEN"

  status_code=$(curl -s -o /dev/null -w "%{http_code}" -X \
    GET $API_URL/v1/marketplace/instances/supplied \
    -H "Authorization: Bearer $TOKEN")

  echo "status code: $status_code"

  if [ $status_code -ne 200 ]; then
    echo "received status code $status_code, expected 200"
    echo "api key is not valid"
    cancel
  fi

  echo "success! api key is valid"
}

check_installed() {
  if ! command -v $1 >/dev/null 2>&1; then
    return 1
  fi
  return 0
}

install_microk8s() {
  echo "Installing microk8s..."
  sudo snap install microk8s --classic --channel=$MICROK8S_VERSION

  sudo snap refresh --hold microk8s
  echo "----------------------"
}

install_hyperdos_if_not_installed() {
  if grep -q hyperdos <<<"$(sudo env "PATH=$PATH" microk8s helm list --all-namespaces)"; then
    echo "hyperdos appears to be installed already, skipping"
  else
    echo "----------------------"
    echo "hyperdos appears not to be installed in the cluster yet, would you like to install it now?"
    if confirm; then
      sudo env "PATH=$PATH" microk8s helm repo add hyperdos https://hyperboliclabs.github.io/Hyper-dOS
      sudo env "PATH=$PATH" microk8s helm install hyperdos hyperdos/hyperdos \
        --version $HYPERDOS_VERSION \
        --set token=$TOKEN $EXTRA_PARAMS
    else
      echo "hyperdos installation canceled by user"
      cancel
    fi
  fi
}

install_microceph() {
  echo "Installing microceph..."
  sudo snap install microceph

  # hard lesson: sometimes canonical will break your entire market without warning
  sudo snap refresh --hold microceph
  echo "----------------------"
}

configure_microceph() {
  # https://docs.ceph.com/en/reef/
  # this install script is designed to set up a single-node cluster
  # so we set the replication factor to 1

  sudo env "PATH=$PATH" microceph.ceph config set global osd_pool_default_size 1
  sudo env "PATH=$PATH" microceph.ceph config set mgr mgr_standby_modules false
  sudo env "PATH=$PATH" microceph.ceph config set osd osd_crush_chooseleaf_type 0
  # modprobe rbd
}

allocate_microceph_disk() {
  # check how much free space is present in this filesystem (.)
  free_space=$(df -kh . | grep '/' | awk '{print $4}')

  # this fills the 'disk_size_gb' variable
  read_disk_size_gb $free_space

  echo "Allocating microceph virtual disk with size: ${disk_size_gb}G"
  # microceph disk add loop,<size in G>,<replication factor>
  sudo env "PATH=$PATH" microceph disk add loop,${disk_size_gb}G,1

  # save 20% of the disk to avoid ceph weirdness
  quota_size_gb=$(((disk_size_gb * 80) / 100))

  # create a resource quota in the instance namespace
  cat <<EOF | sudo env "PATH=$PATH" microk8s kubectl apply -f -
  apiVersion: v1
  kind: ResourceQuota
  metadata:
    name: hyperstore
    namespace: instance
  spec:
    hard:
      # persistentvolumeclaims: "100"  # maximum 100 PVCs
      requests.storage: "${quota_size_gb}G"
EOF

}

# Ask the user how much space they want to allocate to microceph
# and write the value to the disk_size_gb variable
read_disk_size_gb() {
  free_space=$1

  while true; do
    read -r -p "
    Please enter an integer to set the size of the new microceph virtual disk 
    (estimated free space: $free_space)
    Note: it is recommended to leave at least 100GB of free space for ephemeral storage etc.
      
      Enter the number of GB to allocate to the microceph virtual disk: " disk_size_gb

    if [[ $disk_size_gb =~ ^[0-9]+$ ]]; then
      # if (( $disk_size_gb > $free_space )); then
      #   echo "The size specified ($disk_size_gb) is too large. It must be less than the free space on the filesystem."
      break
    else
      echo "Invalid input. Please enter an integer less than $free_space."
    fi
  done
}

cancel() {
  echo "----------------------"
  echo "Installation cancelled"
  echo "----------------------"
  exit
}

confirm() {
  while true; do
    read -r -p "$1 [y/n]: " yn
    case $yn in
    [Yy]*) return 0 ;;
    [Nn]*) return 1 ;;
    *) echo "Please answer yes or no." ;;
    esac
  done
}

count_microk8s_nodes() {
  sudo env "PATH=$PATH" microk8s kubectl get nodes --no-headers | wc -l
}

count_microceph_nodes() {
  sudo env "PATH=$PATH" microceph cluster list | grep ONLINE | wc -l
}

###
## main script
###
echo "----------------------"
echo "Beginning HyperdOS installation..."
echo "----------------------"

check_for_linux

# first, decide whether to install microk8s
if ! check_installed microk8s; then
  echo "microk8s is not installed, would you like to install it now?"

  if confirm; then
    install_microk8s
  else
    cancel
  fi

else
  echo "microk8s is already installed, skipping"
fi

# then, decide whether to install microceph
if ! check_installed microceph; then
  echo "microceph is not installed, would you like to install it now?"

  if confirm; then
    install_microceph
  else
    cancel
  fi

else
  echo "microceph is already installed, skipping"
fi
echo "----------------------"

echo "Starting microk8s..."
sudo env "PATH=$PATH" microk8s start
sudo env "PATH=$PATH" microk8s status --wait-ready
echo "----------------------"

# check if number of nodes is greater than 1
if ((count_microk8s_nodes > 1)); then
  echo "ERROR: microk8s has more than 1 node, this is not currently supported by the install script"
  cancel
fi

echo "----------------------"
echo "Creating namespaces..."
# note: these are idempotent
sudo env "PATH=$PATH" microk8s kubectl create namespace hyperdos || true
sudo env "PATH=$PATH" microk8s kubectl create namespace hyperweb || true
# note that the instance namespace must be created before the hyperstore resourcequota
sudo env "PATH=$PATH" microk8s kubectl create namespace instance || true
sudo env "PATH=$PATH" microk8s kubectl create namespace ping || true
echo "done!"

microceph_node_count=$(count_microceph_nodes)
echo "microceph nodes: $microceph_node_count"

if (($microceph_node_count > 1)); then
  echo "ERROR: microceph has more than 1 node, this is not currently supported by the install script"
  cancel
fi

if (($microceph_node_count == 1)); then
  echo "microceph server appears to be set up already, skipping"
else
  echo "Setting up microceph..."
  # https://microk8s.io/docs/how-to-ceph
  # https://canonical-microceph.readthedocs-hosted.com/en/reef-stable/tutorial/single-node/
  sudo env "PATH=$PATH" microceph cluster bootstrap
  configure_microceph
  echo "done!"
fi
echo "----------------------"

microceph_osd_count="$(sudo env "PATH=$PATH" microceph.ceph osd ls | wc -l)"
echo "microceph disks/osds: $microceph_osd_count"

if (($microceph_osd_count >= 1)); then
  echo "microceph virtual disk appears to be set up already, skipping"
else
  echo "setting up the microceph virtual disk..."
  allocate_microceph_disk
  echo "done!"
fi

echo "----------------------"
echo "microceph.ceph status:"
sudo env "PATH=$PATH" microceph.ceph status

echo "----------------------"
echo "Enabling microk8s components..."
# this is idempotent already
echo "-------------"
sudo env "PATH=$PATH" microk8s enable rbac
echo "-------------"
sudo env "PATH=$PATH" microk8s enable community
echo "-------------"
sudo env "PATH=$PATH" microk8s enable argocd
echo "-------------"
sudo env "PATH=$PATH" microk8s enable nvidia
echo "-------------"
sudo env "PATH=$PATH" microk8s enable rook-ceph

echo "----------------------"
echo "Connecting microk8s to microceph..."
# sudo microk8s connect-external-ceph --no-rbd-pool-auto-create
# TODO check if it's already connected? this is idempotent already though
sudo env "PATH=$PATH" microk8s connect-external-ceph
echo "done!"

echo "----------------------"
namespace="argocd"
while true; do
  pods=$(sudo env "PATH=$PATH" microk8s kubectl get pods -n "$namespace" --no-headers 2>&1)

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

echo "----------------------"
echo "----------------------"
echo "And finally: Installing hyperdos into the cluster..."

# get the token from the user if necessary
check_for_token

# make a test query to the Hyperbolic API
validate_token

install_hyperdos_if_not_installed

echo "==========================="
echo "Installation complete!"
echo "you can view your new cluster at https://app.hyperbolic.xyz/supply"

echo "

@@@@@:@@@:@@@@-:@@@@:-=@@@:#@:@:@:@@@@@@:@@@@@@@@:@@@@:@@:@@@@@:@@@@@@
@@@#=@=-=.@:-:#@#@@@:.@@@:+#*#:@::@*@=:%@-::@@=#*=:@@#@*=@@@=@@:=*#=@@
@=@-=*@=@-::@@+*@=#=@@=%-@::@:@@=#:@*%=-:@-::=@@@@=#:@@@@@#:%::#--#+:@
@*=%#+@::*=#@-+@.%:@*@:=*:-@  =-*: :  .:.@::@: . *@#+#@=@--#@@@@**@=@@
@@@:@-@=@@@-: @-#*     :+  #                     :%:@:@*@=#.::@@=@:@@@
@:+#:@*@@=: .:+ :-::@.:@#+* *: --       *@: #-.:.:.* .. :@:=@=-@#-*.@@
@%:.-@===*+@::@@-*.@::::@-@*..-@@:    :@@:@#@.#:@.::=@::-@:  @ @+:@=@@
@@-=*@:-=@:@@@-@-:=#.#-@:---@.%@:    :@=-@=::=:@@@@.@@@@@@@*=  ::@@+@@
@@=@@+  #::@@@@@#@::@+::.:%@@@@@@    @+@@@@%:=.:@-#@@@@@@@@%#  :=-@@:@
@:+=:   :@@@@@@@@:@@==::*@@@@@@:      -@@@@@@-=-::=@@@@@@@@@#  .-@@:%@
@@@##=:     .#@@@@#=::==@@@-+*@          @:@.@--=@@:@@@@@       @@:@=@
@.:@:@-       =*+@:@==.+@@@--#-%         .-@@#%+*@:@-@:%@        @:%.@
@=::#%:        @@@@=#@=+@@@:#+ =          @:::-:@@%@@@#          #:+:@
@@*#:      ::- :::@*==@=@@:@           :@#@-.@#:%@@@@#-@       @%@@*@@
@-:*=:@     @@=-=:@@%.=:*-.@            ::@@:++=%@@@@*:         %:@:@@
@:@-@+         @+@@#:-=#@#:##   . %     :-@@:+=--#@@#+:         #=+#:@
@@@+%          @@@%*=:@@##@-:#:-.@-:#-: +@@@@:-@@#@@@:%          -@-@@
@@@          +# ::@@==@%+-@-@@--:-@-%-=-::%@:-@=@*@@@@* :         @:@@
@::@#        @#@=@@%#=@@--@@-.:@++@-*..:@::@:##-@-@@+@#@         ::#@@
@@%=%@        ::@%@#@#:@@@#@@@@@@@@@@@@@@@@#@@@+@@@@+=@@@        :#+:@
@@==@:     @@   #@@@@+=#@@@@@@@@@@@@@@@@@@@@@@@-.:@@@:@          @:=:@
@@*:@      @*@#--@@+@%@#@@@@@ %@@@@ @#@@@ #@%@#=%@@@@+:        :#.-@@@
@:#@         +:*%@@@=*%#+++     @     @   #%@@%@@@@@%@:        #*:#:-@
@@:.@@:      :+:@+@@@##@@@@*             @=%.@@@%-@@@=@%:       -=@@=@
@::@==@     :+ :@@@#@@@@@*%-@           :@::-@@##@@@@*:.+*@     .::::@
@:=@:@=:=      .@@@@@#.@@@@=             +@#@@%@+*@+@@.      :#::@=:@@
@@::@#.       -:-%@#@%@@@:*               @=:-#@##@@@@#:-     =@:##=#@
@@-:-+#-@: == @::@@@:@%@@*:=:             ==@#%+%@@@@+@:*   @-=:@@:.@@
@-=:@::@%@@ -+##+@@@+@@@@@=:  =        @+::@@+@+@:@@@@#-    :#:-@.@+@@
@@:=@::@=..%==@@-##@@@@@@@@@:#-        @@::=*@%@@@@#@@##=@@:@:@=@@@:@@
@:@%:@*+-@:+@@*+##@@@@@@%%=@@@@%+*   #-@+@@=:@@@%@#@@@%##@@* :@@@-:+:@
@@::@@.*:::@@@@*@@@@@@+%@:@@@@@@+  -- @@@@@@@@*@@:@@@@@@@@@@@#.:=@@@@@
@:@@*+:@-:*@#@@@@@@@@@@@@@@@@@@@.  =##-#@@@@@@@@@@@@@@@@@@*=.:@==:-:@@
@-#*@=--@@@@%@@@@@@@-#@@@@@@@@@@-: @# :%@@@@@@@%@@@@@@@@@@=@-@-#=#+=@@
@:*@##=#+@@@@@:@:%:@-@@-+@@:: @@@:%%:@@ :%@..#@:%@@@@@@%@:@@=@#+*+@@%@
@@:*%##@#@:*:#::@:@#*@@=@:@%:%:::@@:@*#+.@@@@:###=*-:@*##@.:%#@@  #=:@
@@=*##-=@##@@++:#=@+##=@=-=%@+@@@:=@@:#*@@=:=:*###@.@:@@%:@=@#:@##@@@@
@@@@@+@@:@@@@@@@@@@@@@:=@@@@@@@@@@@@@:@@@@@@@@:@@@@@@@@%@@@@@:@@:@@@@@

"

echo "Welcome to the rAInforest!"
echo "==========================="
