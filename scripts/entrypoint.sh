#!/bin/bash

start_ssh_server() {
  echo "nameserver 1.1.1.1" >>/etc/resolv.conf

  # # this is usually handled by the init container
  # cp -a /home/backup/. /home/ubuntu/
  # cp /tmp/auth_key /home/ubuntu/.ssh/authorized_keys
  # chmod 600 /home/ubuntu/.ssh/authorized_keys
  # chown -R ubuntu: /home/ubuntu

  service ssh start #&& tail -f /dev/null
}

start_ssh_server

pid="$(pgrep -f sshd)"

# Trap the SIGTERM signal and forward it to the sshd process
trap 'kill -SIGTERM $pid; wait $pid' SIGTERM

# Wait for the sshd process to complete
wait "$pid"
