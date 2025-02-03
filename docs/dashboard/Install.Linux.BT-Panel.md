# Linux 宝塔面板安装教程

## 准备工作

预装软件

- Nginx 1.24+ (部分老版本有 `gRPC` 兼容性问题)

## 安装

1. 执行安装命令

```shell
bash <(curl -sSL "https://gitlab.com/CoiaPrant/Sleepy/-/raw/master/scripts/install_dashboard.sh")
```

---

2. 进入 `/opt/sleepy-dashboard` 文件夹

---

3. 创建一个名为 `config.json` 的配置文件

```json
{
    "Web": {
        "Type": "unix",
        "Listen": "/run/sleepy-dashboard.sock"
    },
    "Security": {
        "AllowCORS": false,
        "Secret": "PUT_YOUR_AGENT_SECRET_HERE"
    },
    "System": {
        "DebugMode": false
    }
}
```

> 文件编码必须为 `UTF-8`

`/run/sleepy-dashboard.sock` 为 `Unix Socket` 监听地址, 如果有多个面板请将 `sleepy-dashboard` 更换成不一样的名字

> 请将 `PUT_YOUR_AGENT_SECRET_HERE` 替换为自己的密钥, 建议大小写数字混合

---

4. 修改 `custom/site_config.json`

```json
{
    "title": "Sleepy",
    "description": "",
    "nickname": "TA"
}
```

- `title` 为站点标题
- `description` 为 首页下方显示的 HTML 代码, 可留空
- `nickname` 网站所有者的昵称

---

5. 启动面板

```shell
systemctl enable --now sleepy-dashboard # sleepy-dashboard 为默认服务名, 如果您安装了多个后端请自行修改服务名称
```

6.  在 宝塔面板 (aaPanel) 添加网站并配置 SSL (开启强制 HTTPS)

---

7. 配置反向代理

目标 URL 填写 `http://unix:/run/sleepy-dashboard.sock` (上文的 `Web` 节 `Listen` 字段)

发送域名 `$host`

> 部分版本的面板无法添加, 可以将 目标 URL 设置为 http://127.0.0.1 然后修改反向代理配置文件

编辑反向代理配置文件为以下内容

```nginx
#PROXY-START/
underscores_in_headers on;
 
location ^~ / {
    proxy_pass http://unix:/run/sleepy-dashboard.sock;
    proxy_set_header Host $host;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $http_connection;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header REMOTE-HOST $remote_addr;
 
    add_header X-Cache $upstream_cache_status;
}
 
location ^~ /proto.Sleepy {
    grpc_pass grpc://unix:/run/sleepy-dashboard.sock;
    grpc_read_timeout 300d;
    grpc_send_timeout 300d;
    grpc_socket_keepalive on;
    grpc_set_header Host $host;
    grpc_set_header X-Real-IP $remote_addr;
    grpc_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    grpc_set_header REMOTE-HOST $remote_addr;
    access_log off;
}
 
#PROXY-END/
```

---

8.  使用 **HTTPS** 配合 主机 IP 或 绑定域名 访问

## 搭配 Cloudflare CDN 使用

1. 在 Cloudflare 申请一个回源证书 \([申请教程](https://www.google.com/search?q=cloudflare%E7%94%B3%E8%AF%B715%E5%B9%B4%E5%9B%9E%E6%BA%90%E8%AF%81%E4%B9%A6)\), 并将证书的内容保存至 宝塔面板 (aaPanel) 网站证书中, 私钥内容保存至 宝塔面板 (aaPanel) 网站私钥中

---

2. 在 Cloudflare 后台的 `网络` 选项卡中开启 `gRPC` 和 `WebSocket`

---

3. 在 Cloudflare 后台 的 `SSL/TLS` 选项卡, 修改加密模式为 `完全`
