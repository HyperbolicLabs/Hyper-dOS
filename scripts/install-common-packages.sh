#!/bin/bash

# neofetch is required for the welcome message (motd)
# psmisc is required for the 'killall' command

apt-get update &&
    apt-get install -y \
        openssh-server \
        sudo \
        git \
        neofetch \
        vim \
        psmisc \
        hashcat &&
    apt-get clean

pip3 install jupyterlab
