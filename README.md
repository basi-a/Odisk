# Odisk
不过是一个网盘罢了, 为啥叫odisk, 因为用的对象存储是minio
# TODO
- [x] 脚手架
- [x] GORM
- [X] GIN
- [X] MinIO
- [ ] 数据库设计
- [ ] VUE
- [ ] 前端开发
- [ ] 前后端联调
- [X] 高可用集群模拟
- [ ] 论文
# 网络地址规划
整个网络都是由docker模拟的，以下是`hosts`
```yml
x-hosts-common: &hosts-common
  extra_hosts:
      - "mariadb:172.40.20.10"
      - "adminer:172.40.20.11"
      - "redis:172.40.20.12"
      - "nginx:172.40.20.13"
      - "keepalived_haproxy1:172.40.20.14"
      - "keepalived_haproxy2:172.40.20.15"
      - "minio1:172.40.20.16"
      - "minio2:172.40.20.17"
      - "minio3:172.40.20.18"
      - "minio4:172.40.20.19"
      - "server1:172.40.20.20"
      - "server2:172.40.20.21"
      - "server3:172.40.20.22"
      - "server4:172.40.20.23"
      - "keepalived_vip:172.40.20.100"
```
# 接口文档
# 食用方式
```bash
git clone https://github.com/basi-a/odisk
cd odisk
go mod tidy
go build
cp odisk dev-containers
cd dev-containers 
docker-compose up -d
curl -f http://172.40.20.100:7000/ping
curl -f http://172.40.20.100:9000/minio/health/live
```