# `pog`, a Proxy for Oblivious HTTP by Guardian Project

Implements a simple OHTTP relay per [RFC 9458](https://www.rfc-editor.org/rfc/rfc9458.html).  This is part of a set of client, relay, and gateway developed by Guardian project described more fully in the `ohttp-gp` [repository](https://github.com/johnhess/ohttp-gp).

This forwards requests from an OHTTP client to a pre-configured OHTTP gateway (see `gatewayURL` in `main.go`).

## Building and running

Install prerequisites (golang):

```
curl -O https://dl.google.com/go/go1.22.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile
```

Then run the server on :8080.

```
go run main.go
```

## Sample nginx config

In practice, you'll want to run this behind nginx.  Here's a sample config that assumes you've got certificates provisioned via certbot.

```
server {
    listen 80;
    server_name ohttp-relay.jthess.com;

    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name ohttp-relay.jthess.com;

    ssl_certificate /etc/letsencrypt/live/ohttp-relay.jthess.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/ohttp-relay.jthess.com/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```