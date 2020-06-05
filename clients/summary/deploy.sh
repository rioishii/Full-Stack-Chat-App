docker build -t rioishii/summary .

docker push rioishii/summary

ssh ec2-user@ec2-3-19-55-61.us-east-2.compute.amazonaws.com < update.sh
