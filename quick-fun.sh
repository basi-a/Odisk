#!/bin/bash
go build
cp odisk dev-containers/server
cp -r cert dev-containers/server
cp -r cert dev-containers/haproxy-keepalived
cp -r cert dev-containers/nginx
cd dev-containers || exit
docker compose up -d
cd ..
sleep 30
curl -k -I https://172.40.20.100:7000/ping
echo "View the APP font with browser. https://localhost"