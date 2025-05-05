#!/bin/sh

echo "
This system has been minimized to reduce unneeded packages.
To restore them, you can run the 'unminimize' command.
"

echo "
Note: outside of /home/ubuntu, please be aware that
ephemeral disk usage above 20 GB will trigger a pod reset.
"

# sudo -H -u ubuntu neofetch # sudo isn't always installed
su ubuntu -c "neofetch --source /etc/update-motd.d/hypercloud.ascii"

echo "
...Welcome to Hyperbolic!
"
