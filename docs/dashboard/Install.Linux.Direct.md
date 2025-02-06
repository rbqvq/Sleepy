# Linux 直接安装教程

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
        "Type": "tcp",
        "Listen": ":443",
        "Cert": "cert.pem",
        "Key": "private.key"
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

`:443` 为面板监听地址, 如果有多个面板请将 `443` 更换成不一样的端口, 然后浏览器带端口访问

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

5. 在文件夹下创建 `cert.pem` 和 `private.key`

对应网站证书

> 如果搭配 Cloudflare CDN 使用可使用 回源证书

---

6. 启动面板

```shell
systemctl enable --now sleepy-dashboard # sleepy-dashboard 为默认服务名, 如果您安装了多个后端请自行修改服务名称
```

---

7.  使用 **HTTPS** 配合 主机 IP 或 绑定域名 访问

## 搭配 Cloudflare CDN 使用

1. 在 Cloudflare 申请一个回源证书 \([申请教程](https://www.google.com/search?q=cloudflare%E7%94%B3%E8%AF%B715%E5%B9%B4%E5%9B%9E%E6%BA%90%E8%AF%81%E4%B9%A6)\), 并将证书的内容保存至 宝塔面板 (aaPanel) 网站证书中, 私钥内容保存至 宝塔面板 (aaPanel) 网站私钥中

---

2. 在 Cloudflare 后台的 `网络` 选项卡中开启 `gRPC` 和 `WebSocket`

---

3. 在 Cloudflare 后台 的 `SSL/TLS` 选项卡, 修改加密模式为 `完全`
