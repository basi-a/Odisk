#!/bin/bash
go build
cp odisk dev-containers/server
cd dev-containers || exit
docker comspoe up -d
cd ..
curl -I http://172.40.20.100:7000/ping