#!/bin/bash
go build
cp odisk dev-containers/server
cd dev-containers || exit
docker compose up -d
cd ..
sleep 30
curl -I http://172.40.20.100:7000/ping