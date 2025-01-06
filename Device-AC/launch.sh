#sudo docker build -t device-ac .
#docker run --rm -p 9996:9996 device-ac
#docker run --rm -p 9996:9996 --add-host=host.docker.internal:host-gateway device-ac

if [ "$1" = "local" ]; then
  echo "Running in local mode..."
  sudo docker build -t device-ac .
  docker run --rm -p 9996:9996 device-ac
elif [ "$1" = "remote" ]; then
  echo "Running in remote mode..."
  sudo docker build -t device-ac .
  docker run --rm -p 9996:9996 -p 9999:9999 -e HOST_IP=$(hostname -I | awk '{print $1}') device-ac
else
  echo "Usage: $0 [local|remote]"
  exit 1
fi