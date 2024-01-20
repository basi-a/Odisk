#!/bin/bash
docker compose down -v
docker rmi $(docker images | grep "dev-containers" | awk '{print $1}')
sudo rm -rf db
sudo rm -rf minio
sudo rm -rf static
sudo rm -rf log
sudo rm server/odisk
