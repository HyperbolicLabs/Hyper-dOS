#!/bin/bash

start_ssh_server() {
  echo "nameserver 1.1.1.1" >>/etc/resolv.conf
  service ssh start && tail -f /dev/null
}

start_ssh_server
