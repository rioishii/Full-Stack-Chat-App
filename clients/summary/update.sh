docker rm -f summarytest

docker pull rioishii/summary

docker run \
    -d \
    -p 80:80 \
    -p 443:443 \
    --name summarytest \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    rioishii/summary 

exit