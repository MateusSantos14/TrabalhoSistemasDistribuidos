#docker build -t my-gateway .
#docker run --rm -p 9991:9991 my-gateway
#docker run --rm -p 9991:9991 -p 9996:9996 -p 9997:9997 -p 9998:9998 --add-host=host.docker.internal:host-gateway my-gateway


if [ "$1" = "local" ]; then
  echo "Running in local mode..."
  docker build -t my-gateway .
  docker run --rm -p 9991:9991 my-gateway
elif [ "$1" = "remote" ]; then
  echo "Running in remote mode..."
  docker build -t my-gateway .
  docker run --rm -p 9991:9991 -p 9996:9996 -p 9997:9997 -p 9998:9998 -e HOST_IP=$(hostname -I | awk '{print $1}') my-gateway
else
  echo "Usage: $0 [local|remote]"
  exit 1
fi
