#docker build -t device-carloc .
#docker run --rm -p 9997:9997 device-carloc
#docker run --rm -p 9997:9997 --add-host=host.docker.internal:host-gateway device-carloc

if [ "$1" = "local" ]; then
  echo "Running in local mode..."
  docker build -t device-carloc .
  docker run --rm -p 9997:9997 device-carloc
elif [ "$1" = "remote" ]; then
  echo "Running in remote mode..."
  docker build -t device-carloc .
  docker run --rm -p 9997:9997 -p 9999:9999 -e HOST_IP=$(hostname -I | awk '{print $1}') device-carloc
else
  echo "Usage: $0 [local|remote]"
  exit 1
fi