FROM pytorch/pytorch:2.4.1-cuda12.4-cudnn9-devel
# FROM ubuntu:latest

# Install SSH server
RUN apt-get update && apt-get install -y \
    openssh-server \
    sudo \
    git \
    neofetch \
    vim \
    && apt-get clean

RUN useradd -m -s /bin/bash ubuntu
RUN update-alternatives --install /usr/bin/python python /opt/conda/bin/python 1

RUN adduser ubuntu sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN mkdir /home/ubuntu/.ssh
RUN echo "TODO" > /home/ubuntu/.ssh/authorized_keys

# commented out: this is handled by the entrypoint and volume mount
# RUN chown -R ubuntu:ubuntu /home/ubuntu/.ssh
# RUN chmod 600 /home/ubuntu/.ssh/authorized_keys

# Add CUDA binaries to PATH and CUDA libraries to LD_LIBRARY_PATH
RUN echo 'export PATH="/usr/local/cuda/bin:${PATH}"' >> /etc/bash.bashrc
RUN echo 'export LD_LIBRARY_PATH="/usr/local/cuda/lib64:${LD_LIBRARY_PATH}"' >> /etc/bash.bashrc

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY scripts/99-hypercloud-motd.sh /etc/update-motd.d/99-hypercloud-motd.sh
RUN chmod +x /etc/update-motd.d/99-hypercloud-motd

# switch port to 2222
RUN sed -i 's/#Port 22/Port 2222/' /etc/ssh/sshd_config

# set up a home dir backup to write to the pv later
RUN cp -r /home/ubuntu /home/backup

# SSH port
EXPOSE 22
EXPOSE 2222

ENTRYPOINT ["/entrypoint.sh"]
