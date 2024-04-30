#!/bin/bash
P12_FILE_NAME="server.p12"
PRIVATE_KEY_FILE_NAME="server.key"
CRT_FILE_NAME="server.crt"
PEM_FILE_NAME="server.pem" 
mkcert -cert-file $CRT_FILE_NAME -key-file $PRIVATE_KEY_FILE_NAME \
    "*.basi-a.top" localhost 127.0.0.1 ::1 172.40.20.1 172.40.20.100 172.40.20.13 172.40.20.20 172.40.20.21 172.40.20.22 172.40.20.23 minio1 minio2 minio3 minio4
openssl pkcs12 -export -inkey $PRIVATE_KEY_FILE_NAME -in $CRT_FILE_NAME -out $P12_FILE_NAME
#password basi1024
openssl x509 -in $CRT_FILE_NAME -text -noout
cat $PRIVATE_KEY_FILE_NAME $CRT_FILE_NAME > $PEM_FILE_NAME
mkcert -install


# 这样弄出来的，浏览器就不会报不安全红锁了，而且系统信任 curl 不用加-k了
