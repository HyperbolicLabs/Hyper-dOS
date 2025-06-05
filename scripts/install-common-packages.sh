#!/bin/bash

# neofetch is required for the welcome message (motd)
# psmisc is required for the 'killall' command

set -e

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
                rsync && \
        apt-get clean

curl -s https://ollama.com/install.sh | bash

# pip3 install jupyterlab
su ubuntu -c "pip3 install jupyterlab"
echo 'export PATH="/home/ubuntu/.local/bin:$PATH"' >> /home/ubuntu/.profile

