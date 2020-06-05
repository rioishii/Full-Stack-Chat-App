docker rm -f mongodb

docker run -d \
    --name mongodb \
    -p 27017:27017 \
    mongo