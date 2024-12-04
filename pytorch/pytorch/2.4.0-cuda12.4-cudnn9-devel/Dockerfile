FROM pytorch/pytorch:2.4.0-cuda12.4-cudnn9-devel

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

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# switch port to 2222
RUN sed -i 's/#Port 22/Port 2222/' /etc/ssh/sshd_config

# SSH port
EXPOSE 22
EXPOSE 2222

# USER ubuntu

# chown authorized keys file and start the SSH service. Then keep the container alive
ENTRYPOINT ["/entrypoint.sh"]
