#!/bin/bash
docker compose down -v
sudo rm -rf db
sudo rm -rf minio
sudo rm -rf static
sudo rm -rf log
sudo rm server/odisk
