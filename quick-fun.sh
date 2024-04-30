#!/bin/bash
go build
cp odisk dev-containers/server
cp -r template dev-containers/server
cp -r cert dev-containers/server
cp -r cert dev-containers/haproxy-keepalived
cp -r cert dev-containers/nginx
cp -r cert dev-containers/minio
cd dev-containers || exit
docker compose up -d
cd ..
sleep 30
curl -k -I https://localhost/api/ping
echo "View the APP font with browser. https://localhost"
