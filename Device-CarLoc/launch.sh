DEVICE_ID=${1:-"CL-1"}
PORT=${2:-9997}

sudo cp ../DeviceClasses/SimulatedSensor.py SimulatedSensor.py
sudo docker build -t device-carloc .
#docker run -p 9998:9998 --network my-network device-carloc
sudo docker run --rm -p "$PORT:$PORT" device-carloc "$DEVICE_ID" "$PORT"
