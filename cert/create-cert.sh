#!/bin/bash
P12_FILE_NAME="server.p12"
PRIVATE_KEY_FILE_NAME="server.key"
CRT_FILE_NAME="server.crt"

mkcert -cert-file $CRT_FILE_NAME -key-file $PRIVATE_KEY_FILE_NAME \
    "*.basi-a.top" localhost 127.0.0.1 ::1 172.40.20.1 172.40.20.100 172.40.20.13 172.40.20.20 172.40.20.21 172.40.20.22 172.40.20.23
openssl pkcs12 -export -inkey $PRIVATE_KEY_FILE_NAME -in $CRT_FILE_NAME -out $P12_FILE_NAME
#password basi1024
openssl x509 -in $CRT_FILE_NAME -text -noout

mkcert -install