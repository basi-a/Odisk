#!/bin/bash
# docker compose down -v
# docker rmi $(docker images | grep "odisk" | awk '{print $1":latest"}')
# # sudo rm -rf db
# # sudo rm -rf minio   # 注释掉是为了不清理掉已经注册的用户、存储桶、以及二者之间的关系映射
# sudo rm -rf static
# sudo rm -rf log
# sudo rm -rf nsq
# sudo rm server/odisk
# sudo rm -rf server/cert
# sudo rm -rf server/template
# sudo rm -rf haproxy-keepalived/cert
# sudo rm -rf nginx/cert

# 检查是否有清理数据的参数
CLEAN_DATA=$1
if [ "$CLEAN_DATA" = "--clean-data" ]; then
    echo "Cleaning MinIO and Database data..."
    sudo rm -rf db
    sudo rm -rf minio

fi

docker compose down -v
docker rmi $(docker images | grep "odisk" | awk '{print $1":latest"}')
# 不论是否提供清理数据参数，都会执行的清理操作
sudo rm -rf static
sudo rm -rf log
sudo rm -rf nsq
sudo rm server/odisk
sudo rm -rf server/cert
sudo rm -rf server/template
sudo rm -rf haproxy-keepalived/cert
sudo rm -rf nginx/cert
