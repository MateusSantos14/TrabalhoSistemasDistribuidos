DEVICE_ID=${1:-"HL-1"}
PORT=${2:-9998}

sudo cp ../DeviceClasses/SimulatedActuator.py SimulatedActuator.py
sudo docker build -t device-headlight .
#docker run -p 9998:9998 --network my-network device-headlight
sudo docker run --rm -p "$PORT:$PORT" device-headlight "$DEVICE_ID" "$PORT"
