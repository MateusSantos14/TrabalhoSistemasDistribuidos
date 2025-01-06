docker build -t my-gateway .
#docker run -p 9990:9990 --network my-network my-gateway
docker run --rm -p 9990:9990 my-gateway
