cp ../DeviceClasses/SimulatedActuator.py SimulatedActuator.py
docker build -t device-headlight .
#docker run -p 9998:9998 --network my-network device-headlight
docker run -p 9998:9998 device-headlight
