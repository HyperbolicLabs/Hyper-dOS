FROM nvidia/cuda:12.4.1-cudnn-devel-ubuntu22.04 

# RUN adduser --disabled-password --gecos '' docker
RUN useradd -m -s /bin/bash ubuntu

# RUN adduser ubuntu sudo
# RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN mkdir /home/ubuntu/.ssh
RUN echo "TODO" > /home/ubuntu/.ssh/authorized_keys

COPY scripts/install-common-packages.sh /install-common-packages.sh
RUN chmod +x /install-common-packages.sh
RUN /install-common-packages.sh

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

ENTRYPOINT ["/entrypoint.sh"]
