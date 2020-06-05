npm run build

docker build -t rioishii/messages2 .

docker push rioishii/messages2

ssh ec2-user@ec2-3-19-55-61.us-east-2.compute.amazonaws.com < update.sh