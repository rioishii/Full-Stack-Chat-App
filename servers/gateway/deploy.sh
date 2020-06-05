GOOS=linux go build
docker build -t rioishii/gateway .
go clean
docker push rioishii/gateway

cd ../db
docker build -t rioishii/db .
docker push rioishii/db
cd -

cd ../summary 
GOOS=linux go build
docker build -t rioishii/summary .
go clean
docker push rioishii/summary
cd -

cd ../messaging
docker build -t rioishii/messaging .
docker push rioishii/messaging
cd -

ssh ec2-user@ec2-3-20-204-157.us-east-2.compute.amazonaws.com < update.sh