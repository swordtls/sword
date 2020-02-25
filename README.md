# sword proxy
 TLS double check proxy 
##

# run mode
- server

- proxy_server
 
- client
 
- base64p12 (generate p12 perm to ios shadowrocket https cert base64 code)

# server and client mode

source -----> client(mode) ---------> server(mode) ---> remote

# proxy server and client mode

ss local ---->client(mode) -------->proxy_server(mode) ----> ss server

# generate ca and cert

openssl req -newkey rsa:2048 -nodes -keyout ca.key -x509 -days 1024 -out ca.crt -subj "/C=US/ST=CL/L=CF/O=sword.tls/OU=dev/CN=tls_ca/emailAddress=ca@sword.tls"

openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 1024 -out server.pem -subj "/C=US/ST=CL/L=CF/O=sword.tls/OU=dev/CN=tls_server/emailAddress=server@sword.tls"

openssl req -newkey rsa:2048 -nodes -keyout client.key -x509 -days 1024 -out client.pem -subj "/C=US/ST=CL/L=CF/O=sword.tls/OU=dev/CN=tls_client/emailAddress=client@sword.tls"

openssl pkcs12 -export -in client.pem -inkey client.key -out client.p12 -certfile ca.crt

# install 

git clone https://github.com/swordtls/sword.git

cd sword 


*must be run on root user

chmod + x install.sh

./install.sh
