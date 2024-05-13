#!/bin/bash
CONTAINERS_NAME=(
    odisk-server1
    odisk-server2
    odisk-server3
    odisk-server4
    odisk-minio1
    odisk-minio2
    odisk-minio3
    odisk-minio4
    odisk-keepalived_haproxy1
    odisk-keepalived_haproxy2
    odisk-nsqlookupd
    odisk-nsqd
    odisk-nsqadmin
    odisk-nginx
    odisk-postgres
    odisk-adminer
    odisk-redis
)
HEALTH_CHECK_GROUP=(
    odisk-server1
    odisk-server2
    odisk-server3
    odisk-server4
    odisk-minio1
    odisk-minio2
    odisk-minio3
    odisk-minio4
    odisk-nsqlookupd
    odisk-nsqd
    odisk-nginx
    odisk-postgres
    odisk-redis
)

# 检查容器运行状态
check_container_status() {
    container_name=$1
    status=$(docker inspect "$container_name" | jq -r ".[].State.Status")
    format_output "$container_name" "$status"
}

# 检查容器健康状态
check_container_health() {
    container_name=$1
    # 确保容器支持健康检查
    if [[ " ${HEALTH_CHECK_GROUP[*]} " =~ ${container_name} ]]; then
        health_status=$(docker inspect "$container_name" | jq -r ".[].State.Health.Status")
        format_output "$container_name" "$health_status"
    fi
}

# 输出格式化的函数
format_output() {
    printf "  %-30s\t%-12s\n" "$1" "$2"
}

# 添加参数处理逻辑
if [ "$#" -ne 1 ]; then
    echo "  Usage: $0 {status|health}"
    exit 1
fi

check_type=$1

case $check_type in
    status)
        echo "容器运行状态:"
        for container_name in "${CONTAINERS_NAME[@]}"; do
            check_container_status "$container_name"
        done
        ;;
    health)
        echo "容器健康状态:"
        for container_name in "${CONTAINERS_NAME[@]}"; do
            check_container_health "$container_name"
        done
        ;;
    *)
        echo "Invalid option. Please choose 'status' or 'health'."
        exit 1
        ;;
esac
