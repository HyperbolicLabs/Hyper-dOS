FROM pytorch/pytorch:latest

# Install SSH server
RUN apt-get update && \
    apt-get install -y openssh-server sudo && \
    apt-get clean

# RUN adduser --disabled-password --gecos '' docker
RUN useradd -m -s /bin/bash sshuser

RUN adduser sshuser sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN mkdir /home/sshuser/.ssh
RUN echo "TODO" > /home/sshuser/.ssh/authorized_keys

# commented out: this is handled by the entrypoint and volume mount
# RUN chown -R sshuser:sshuser /home/sshuser/.ssh
# RUN chmod 600 /home/sshuser/.ssh/authorized_keys

# Permit root login via SSH
# RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# switch port to 2222
RUN sed -i 's/#Port 22/Port 2222/' /etc/ssh/sshd_config

# SSH port
EXPOSE 22
EXPOSE 2222

# USER sshuser

# chown authorized keys file and start the SSH service. Then keep the container alive
ENTRYPOINT cat /home/sshuser/auth/auth > /home/sshuser/.ssh/authorized_keys && chown -R sshuser:sshuser /home/sshuser/.ssh && chmod 600 /home/sshuser/.ssh/authorized_keys && service ssh start && tail -f /dev/null
