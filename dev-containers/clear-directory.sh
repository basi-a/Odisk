#!/bin/bash
docker compose down -v
docker rmi $(docker images | grep "odisk" | awk '{print $1":latest"}')
# sudo rm -rf db
# sudo rm -rf minio   # 注释掉是为了不清理掉已经注册的用户、存储桶、以及二者之间的关系映射
sudo rm -rf static
sudo rm -rf log
sudo rm -rf nsq
sudo rm server/odisk
sudo rm -rf server/cert
sudo rm -rf server/template
sudo rm -rf haproxy-keepalived/cert
sudo rm -rf nginx/cert