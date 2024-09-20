FROM nvidia/cuda:12.3.1-devel-ubuntu22.04

# Install SSH server
RUN apt-get update && \
    apt-get install -y openssh-server sudo git && \
    apt-get clean

# RUN adduser --disabled-password --gecos '' docker
RUN useradd -m -s /bin/bash ubuntu

RUN adduser ubuntu sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN mkdir /home/ubuntu/.ssh
RUN echo "TODO" > /home/ubuntu/.ssh/authorized_keys

# Add CUDA binaries to PATH and CUDA libraries to LD_LIBRARY_PATH
ENV PATH="/usr/local/cuda/bin:${PATH}"
ENV LD_LIBRARY_PATH="/usr/local/cuda/lib64:${LD_LIBRARY_PATH}"

# switch port to 2222
RUN sed -i 's/#Port 22/Port 2222/' /etc/ssh/sshd_config

# SSH port
EXPOSE 22
EXPOSE 2222

# USER ubuntu

# chown authorized keys file and start the SSH service. Then keep the container alive
ENTRYPOINT cat /home/ubuntu/auth/auth > /home/ubuntu/.ssh/authorized_keys && chown -R ubuntu:ubuntu /home/ubuntu/.ssh && chmod 600 /home/ubuntu/.ssh/authorized_keys && service ssh start && tail -f /dev/null
