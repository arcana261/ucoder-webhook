# ucoder-webhook
A webhook parser for my site!

Install webhook as `https://webhook.ucoder.ir/webhooks/`

# Installation

## Create systemd unit file

```bash
sudo vim /etc/systemd/system/ucoder-webhook.service
```

```
[Unit]
Description=Ucoder Github Webhook Parser
After=network.target

[Service]
Environment="PATH=bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:/home/arcana/.local/bin:/home/arcana/bin"
ExecStart=/usr/local/bin/ucoder-webhook --port 3900 --secret <SOME SECRET>
Restart=on-abort

[Install]
WantedBy=multi-user.target
```
## Deploy

```bash
make deploy
```

## Configure Nginx

```bash
sudo vim /etc/nginx/conf.d/webhook.ucoder.ir.conf
```

```
server {
  listen 80;
  server_name webhook.ucoder.ir;

  access_log /var/log/nginx/access_webhook.log;
  error_log  /var/log/nginx/error_webhook.log;

  location / {
    proxy_pass http://127.0.0.1:3900;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_read_timeout 43200000;

    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header Host $http_host;
    proxy_set_header X-NginX-Proxy true;
  }

}
```

Test nginx configuration

```bash
sudo nginx -t
```

Restart nginx

```bash
sudo systemctl restart nginx
```

allow nginx to connect to webhook

```bash
sudo semanage port -a -t http_port_t -p tcp 3900
sudo /usr/sbin/setsebool -P httpd_can_network_connect 1
```

## Enable SSL

Now install certbot to enable SSL

```bash
yum -y install yum-utils
yum-config-manager --enable rhui-REGION-rhel-server-extras rhui-REGION-rhel-server-optional
sudo yum install certbot-nginx
```

Enable CertBot on our Nginx Configuration
```bash
sudo certbot --nginx
```

**OR IF YOU RUN INTO BUG PROBLEM**
```bash
sudo certbot --authenticator standalone --installer nginx --pre-hook "systemctl stop nginx" --post-hook "systemctl stop nginx"
```

Test if renewal procedure succeeds
```bash
sudo certbot renew --dry-run
```

Enable CronTab job to renew our certificates automatically
```bash
sudo crontab -e
```
```
0 0,12 * * * python -c 'import random; import time; time.sleep(random.random() * 3600)' && certbot renew 
```

## Create DHParam to make it more Secure

```bash
sudo mkdir -p /etc/nginx/ssl/webhook.ucoder.ir
sudo openssl dhparam -out /etc/nginx/ssl/webhook.ucoder.ir/dhparam.pem 4096
```

sudo vim /etc/nginx/conf.d/webhook.ucoder.ir.conf
...............................
ssl_dhparam /etc/nginx/ssl/webhook.ucoder.ir/dhparam.pem;
...............................

**test nginx configuration**
```bash
nginx -t
```

**restart nginx**
```bash
sudo systemctl restart nginx
```
