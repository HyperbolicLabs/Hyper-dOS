#!/bin/bash

# neofetch is required for the welcome message (motd)
# psmisc is required for the 'killall' command

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
                hashcat &&
        apt-get clean

curl -s https://ollama.com/install.sh | bash

# pip3 install jupyterlab
su ubuntu -c "pip3 install jupyterlab"
