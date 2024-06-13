# Odisk
这个是俺滴毕业设计
不过是一个网盘罢了, 为啥叫odisk, 因为用的对象存储是minio
这个要和前端的`odisk-font`放在同一个目录，不然复制ssl证书会有问题
# TODO
- [x] 脚手架
- [x] GORM
- [x] GIN
- [x] MinIO
- [x] Nsq
- [X] 数据库设计
- [x] vue.js
- [x] https
- [x] 前端开发
- [x] 前后端联调
- [x] 高可用集群模拟
- [x] 论文
# 网络地址规划
整个网络都是由docker模拟的，以下是`hosts`
```yml
x-hosts-common: &hosts-common
  extra_hosts:
      - "db:172.30.20.10"
      - "adminer:172.30.20.11"
      - "redis:172.30.20.12"
      - "nginx:172.30.20.13"
      - "keepalived.haproxy1:172.30.20.14"
      - "keepalived.haproxy2:172.30.20.15"
      - "minio1:172.30.20.16"
      - "minio2:172.30.20.17"
      - "minio3:172.30.20.18"
      - "minio4:172.30.20.19"
      - "server1:172.30.20.20"
      - "server2:172.30.20.21"
      - "server3:172.30.20.22"
      - "server4:172.30.20.23"
      - "nsqlookupd:172.30.20.24"
      - "nsqd:172.30.20.25"
      - "nsqadmin:172.30.20.26"
```

# 食用方式
```bash
cd odisk/cert
./create-cert.sh
cd ../ 
cd odisk-font
./build.sh
cd ../
cd odisk
go mod tidy
./quick-fun.sh
```

需要先有`docker`、`docker-compose`、`golang`, 没有的话需要先安装
```bash
sudo pacman -S docker docker-compose go
```

# 生成测试文件
```bash
dd if=/dev/zero of=testfile-a bs=1M count=1024  #1G 正好不用分片的最大文件
dd if=/dev/zero of=testfile-b bs=4M count=1024  #4G 
```