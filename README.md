# Odisk
不过是一个网盘罢了, 为啥叫odisk, 因为用的对象存储是minio
# TODO
- [x] 脚手架
- [x] GORM
- [x] GIN
- [x] MinIO
- [ ] 日志收集
- [X] 数据库设计
- [x] VUE
- [x] https
- [ ] 前端开发
- [ ] 前后端联调
- [x] 高可用集群模拟
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
      - "keepalived.haproxy1:172.40.20.14"
      - "keepalived.haproxy2:172.40.20.15"
      - "minio1:172.40.20.16"
      - "minio2:172.40.20.17"
      - "minio3:172.40.20.18"
      - "minio4:172.40.20.19"
      - "server1:172.40.20.20"
      - "server2:172.40.20.21"
      - "server3:172.40.20.22"
      - "server4:172.40.20.23"
      - "keepalived.vip:172.40.20.100"
```
# 接口文档
[Apifox 文档分享 ](https://apifox.com/apidoc/shared-60f72b42-a39e-4e18-85b5-a0c4e84e415d)
# 食用方式
```bash
git clone https://github.com/basi-a/odisk
git clone https://github.com/basi-a/odisk-font
cd odisk-font
npm install
npm run build
cp -r dist ../odisk/dev-containers/nginx/dist
cd ../odisk
./cert/create-cert.sh 
go mod tidy
./quick-fun.sh
```
需要先有`docker`、`docker-compose`、`golang`, 没有的话需要先安装
```bash
sudo pacman -S docker docker-compose go
```