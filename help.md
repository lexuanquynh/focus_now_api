
service:
```
[Unit]
Description=Focus now APIs pre-production
After=multi-user.target

[Service]
User=root
Group=root
Type=simple
Restart=always
RestartSec=5s
ExecStart=/usr/local/src/focus_now/auth

[Install]
WantedBy=multi-user.target

```
run service:
```
sudo systemctl start focus_now.service
sudo systemctl enable focus_now.service
sudo systemctl status focus_now.service
```

vi /etc/nginx/sites-available/focus_now
```python
server {
        listen 80;

        location /focus_dev {
                proxy_pass http://127.0.0.1:8080/api/v1;
                proxy_set_header Host $host;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }         
}

```
Create a symbolic link of our config file to the sites-enabled folder:
```
ln -s /etc/nginx/sites-available/focus_now /etc/nginx/sites-enabled/focus_now
```

run service:
```
sudo systemctl stop focus_now.service
sudo systemctl start focus_now.service
sudo systemctl enable focus_now.service
sudo systemctl status focus_now.service
```

Tạo https:
```
sudo certbot --nginx -d focus.codetoanbug.com
```