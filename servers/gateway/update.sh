docker rm -f gatewaytest
docker pull rioishii/gateway

docker rm -f mysqlServer
docker pull rioishii/db

docker rm -f redisServer

docker rm -f summary
docker pull rioishii/summary

docker rm -f mongodb
docker rm -f rabbitmq

docker rm -f chat1
docker rm -f chat2
docker rm -f chat3
docker pull rioishii/messaging

export TLSCERT=/etc/letsencrypt/live/api.rioishii.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.rioishii.me/privkey.pem

docker network create mysqlNet

export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)

docker run -d \
    --name mysqlServer \
    --network mysqlNet \
    -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
    -e MYSQL_DATABASE=userDB \
    rioishii/db

docker run -d \
    --name redisServer \
    --network mysqlNet \
    redis

docker run -d \
    --name summary \
    --network mysqlNet \
    -e ADDR=:4000 \
    -p 4000:4000 \
    rioishii/summary

docker run -d \
    -p 27017:27017 \
    --name mongodb \
    --network mysqlNet \
    mongo

docker run -d \
    --hostname myrabbitmq \
    --name rabbitmq \
    -p 5672:5672 -p 15672:15672 \
    --network mysqlNet \
    rabbitmq:3-management 

sleep 10

docker run -d \
    -e ADDR=:5001 \
    --name chat1 \
    --network mysqlNet \
    --restart always \
    rioishii/messaging

docker run -d \
    -e ADDR=:5002 \
    --name chat2 \
    --network mysqlNet \
    --restart always \
    rioishii/messaging

docker run -d \
    -e ADDR=:5003 \
    --name chat3 \
    --network mysqlNet \
    --restart always \
    rioishii/messaging

docker run \
    -d \
    -p 443:443 \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
    -e SUMMARY=summary:4000 \
    -e CHAT="chat1:5001,chat2:5002,chat3:5003" \
    -e RABBITADDR="amqp://rabbitmq:5672" \
    -e TLSCERT=$TLSCERT \
    -e TLSKEY=$TLSKEY \
    --name gatewaytest \
    --network mysqlNet \
    --restart always \
    rioishii/gateway

exit