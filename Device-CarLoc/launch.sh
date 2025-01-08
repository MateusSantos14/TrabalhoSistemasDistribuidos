DEVICE_ID=${1:-"CL-1"}
PORT=${2:-9997}

cp ../DeviceClasses/SimulatedSensor.py SimulatedSensor.py
docker build -t device-carloc .
#docker run -p 9998:9998 --network my-network device-carloc
docker run --rm -p "$PORT:$PORT" device-carloc "$DEVICE_ID" "$PORT"
