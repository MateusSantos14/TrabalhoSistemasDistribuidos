docker build -t device-carloc .
#docker run -p 9998:9998 --network my-network device-carloc
docker run --rm -p 9997:9997 device-carloc
