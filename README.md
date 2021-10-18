# tls-client


### How to build

```
go build
```


### How to generate cert file and key file

```
mkdir certs
rm certs/*

# for tls server (ethpool)
openssl req -new -nodes -x509 -out certs/server.pem -keyout certs/server.key -days 36500

# for tls client (tls-client)
openssl req -new -nodes -x509 -out certs/client.pem -keyout certs/client.key -days 36500
```
