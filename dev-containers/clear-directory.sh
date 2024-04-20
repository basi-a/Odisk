#!/bin/bash
docker compose down -v
docker rmi $(docker images | grep "odisk" | awk '{print $1":latest"}')
sudo rm -rf db
sudo rm -rf minio
sudo rm -rf static
sudo rm -rf log
sudo rm -rf nsq
sudo rm server/odisk
sudo rm -rf server/cert
sudo rm -rf server/template
sudo rm -rf haproxy-keepalived/cert
sudo rm -rf nginx/cert