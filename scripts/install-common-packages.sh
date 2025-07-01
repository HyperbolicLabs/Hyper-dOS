#!/bin/bash

# neofetch is required for the welcome message (motd)
# psmisc is required for the 'killall' command

set -e

export DEBIAN_FRONTEND=noninteractive

apt-get clean

apt-get update &&
        apt-get install -y \
                openssh-server \
                sudo \
                git \
                neofetch \
                curl \
                vim \
                psmisc \
                python3-pip \
                hashcat \
                mosh \
                rsync &&
        apt-get clean

adduser ubuntu sudo
echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers.d/99-hyperdos-sudo
# alternatively, no need to adduser.
# just write "ubuntu ALL=(ALL) NOPASSWD:ALL" - this appears to be how they do it in ec2/cloud-init


curl -s https://ollama.com/install.sh | bash

# pip3 install jupyterlab
echo 'export PATH="/home/ubuntu/.local/bin:$PATH"' >>/home/ubuntu/.profile
su ubuntu -c "pip3 install jupyterlab"
