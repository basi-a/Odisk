#!/bin/bash
go build
cp odisk ./dev-containers
cd dev-containers || exit 
docker-compose up -d
curl -I http://172.40.20.100:7000/ping