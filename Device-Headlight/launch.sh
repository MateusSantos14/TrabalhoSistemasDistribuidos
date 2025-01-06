#docker build -t device-headlight .
#docker run --rm -p 9998:9998 device-headlight
#docker run --rm -p 9998:9998 --add-host=host.docker.internal:host-gateway device-headlight

#!/bin/bash

if [ "$1" = "local" ]; then
  echo "Running in local mode..."
  docker build -t device-headlight .
  docker run --rm -p 9998:9998 device-headlight
elif [ "$1" = "remote" ]; then
  echo "Running in remote mode..."
  docker build -t device-headlight .
  docker run --rm -p 9998:9998 -p 9999:9999 -e HOST_IP=$(hostname -I | awk '{print $1}') device-headlight
else
  echo "Usage: $0 [local|remote]"
  exit 1
fi