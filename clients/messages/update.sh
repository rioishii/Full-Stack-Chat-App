docker rm -f messages2

docker pull rioishii/messages2

docker run \
    -d \
    -p 80:80 \
    -p 443:443 \
    --name messages2 \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    rioishii/messages2

exit