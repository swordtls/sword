#!/bin/bash
yum install wget -y
apt install wget -y

#create key and perm
openssl req -newkey rsa:2048 -nodes -keyout ca.key -x509 -days 1024 -out ca.crt -subj "/C=US/ST=CL/L=CF/O=sword.tls/OU=dev/CN=tls_ca/emailAddress=ca@sword.tls"
openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 1024 -out server.pem -subj "/C=US/ST=CL/L=CF/O=sword.tls/OU=dev/CN=tls_server/emailAddress=server@sword.tls"
openssl req -newkey rsa:2048 -nodes -keyout client.key -x509 -days 1024 -out client.pem -subj "/C=US/ST=CL/L=CF/O=sword.tls/OU=dev/CN=tls_client/emailAddress=client@sword.tls"
openssl pkcs12 -export -in client.pem -inkey client.key -out client.p12 -certfile ca.crt

#get go binary file
wget https://dl.google.com/go/go1.13.8.linux-amd64.tar.gz

tar -C /tmp -zxf go1.13.8.linux-amd64.tar.gz

#build sword
/tmp/go/bin/go build sword.go

./sword --mode=b
./sword --mode=s --conf=server.json --daemon=true 2>&1 >> /tmp/sword.log&