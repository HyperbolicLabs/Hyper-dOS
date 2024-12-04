#!/bin/bash

start_ssh_server() {
  echo "nameserver 1.1.1.1" >>/etc/resolv.conf
  cat /tmp/auth_key >/home/ubuntu/.ssh/authorized_keys
  chown -R ubuntu:ubuntu /home/ubuntu/.ssh
  chmod 600 /home/ubuntu/.ssh/authorized_keys
  service ssh start && tail -f /dev/null
}

start_ssh_server
