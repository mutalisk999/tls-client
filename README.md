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


### Add to service

`sudo vi /etc/systemd/system/tls-client.service`


```
[Unit]
After=network.target

[Service]
User=root
WorkingDirectory=/home/ubuntu/tls-client
ExecStart=/home/ubuntu/tls-client/tls-client -d /home/ubuntu/tls-client
```

### Service start/stop/restart/status

```
sudo systemctl start tls-client.service
sudo systemctl stop tls-client.service
sudo systemctl restart tls-client.service
sudo systemctl status tls-client.service
```


### Enable start service while system bootint up

```
sudo systemctl enable tls-client.service
sudo systemctl is-enabled tls-client.service
```


