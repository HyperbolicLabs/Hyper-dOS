#!/bin/bash

start_ssh_server() {
  echo "nameserver 1.1.1.1" >>/etc/resolv.conf

  # # this is usually handled by the init container
  # cp -a /home/backup/. /home/ubuntu/
  # cp /tmp/auth_key /home/ubuntu/.ssh/authorized_keys
  # chmod 600 /home/ubuntu/.ssh/authorized_keys
  # chown -R ubuntu: /home/ubuntu

  service ssh start && tail -f /dev/null
}

# Kubernetes sends a SIGTERM signal to the main process in a container when terminating a pod
# in our case, we have an sshd process running in the container that needs to be cleaned up
# otherwise, the container will not terminate smoothly
# this can lead to problems with clean eviction and rescheduling of the instance
# Note: these can theoretically be replaced with a 'preStop.exec.command' lifecycle parameter in the pod yaml.
trap 'killall sshd' EXIT # this 'EXIT' appears to be what actually occurs
trap 'killall sshd' TERM # but including 'TERM' for extra coverage

start_ssh_server
