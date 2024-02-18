#!/bin/bash
PRIVATE_KEY_FILE_NAME="server.key"
CRT_FILE_NAME="server.crt"
PEM_FILE_NAME="server.pem"
echo "Create Private Key"
openssl ecparam -genkey -name secp384r1 -out $PRIVATE_KEY_FILE_NAME
echo "Create Cert"
openssl req -new -x509 -sha256 -key $PRIVATE_KEY_FILE_NAME -out $CRT_FILE_NAME -days 3650
echo "Check Out"
openssl x509 -in $CRT_FILE_NAME -text -noout
echo "merge"
cat $CRT_FILE_NAME $PRIVATE_KEY_FILE_NAME > $PEM_FILE_NAME