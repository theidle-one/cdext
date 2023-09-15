#!/bin/bash
if ! [ $(id -u) = 0 ]; then
   echo "The script need to be run as root." >&2
   exit 1
fi

if [ $SUDO_USER ]; then
    real_user=$SUDO_USER
else
    real_user=$(whoami)
fi

export LOG_MGT_DASHBOARD=http://localhost:8081
export LOG_MGT_BACKEND=http://localhost:8080


caddy run -watch