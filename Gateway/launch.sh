sudo docker build -t my-gateway .
#docker run -p 9990:9990 --network my-network my-gateway
sudo docker run --rm -p 9991:9991 my-gateway
