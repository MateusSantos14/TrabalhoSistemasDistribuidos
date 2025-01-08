DEVICE_ID=${1:-"HL-1"}
PORT=${2:-9998}

cp ../DeviceClasses/SimulatedActuator.py SimulatedActuator.py
docker build -t device-headlight .
#docker run -p 9998:9998 --network my-network device-headlight
docker run --rm -p "$PORT:$PORT" device-headlight "$DEVICE_ID" "$PORT"
