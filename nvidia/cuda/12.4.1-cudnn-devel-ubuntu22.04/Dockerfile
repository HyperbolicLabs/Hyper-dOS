FROM nvidia/cuda:12.4.1-cudnn-devel-ubuntu22.04 

# Install SSH server
RUN apt-get update && \
    apt-get install -y \
    openssh-server \
    sudo \
    git \
    neofetch \
    && apt-get clean

# RUN adduser --disabled-password --gecos '' docker
RUN useradd -m -s /bin/bash ubuntu

RUN adduser ubuntu sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN mkdir /home/ubuntu/.ssh
RUN echo "TODO" > /home/ubuntu/.ssh/authorized_keys

# Add CUDA binaries to PATH and CUDA libraries to LD_LIBRARY_PATH
RUN echo 'export PATH="/usr/local/cuda/bin:${PATH}"' >> /etc/bash.bashrc
RUN echo 'export LD_LIBRARY_PATH="/usr/local/cuda/lib64:${LD_LIBRARY_PATH}"' >> /etc/bash.bashrc

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# clean out default motd
RUN rm /etc/update-motd.d/*

# install hypercloud motd script
COPY scripts/99-hypercloud-motd.sh /etc/update-motd.d/99-hypercloud-motd
COPY scripts/hypercloud.ascii /etc/update-motd.d/hypercloud.ascii
RUN chmod +x /etc/update-motd.d/99-hypercloud-motd

# switch port to 2222
RUN sed -i 's/#Port 22/Port 2222/' /etc/ssh/sshd_config

# SSH port
EXPOSE 22
EXPOSE 2222

# USER ubuntu

# chown authorized keys file and start the SSH service. Then keep the container alive
ENTRYPOINT ["/entrypoint.sh"]
