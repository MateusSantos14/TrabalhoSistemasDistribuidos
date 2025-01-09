DEVICE_ID=${1:-"AC-1"}
PORT=${2:-9996}

sudo cp ../DeviceClasses/SimulatedActuator.py SimulatedActuator.py
sudo docker build -t device-ac .
#docker run -p 9996:9996 --network my-network device-ac
sudo docker run --rm -p "$PORT:$PORT" device-ac "$DEVICE_ID" "$PORT"