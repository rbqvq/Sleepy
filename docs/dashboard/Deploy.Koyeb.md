# Koyeb 平台免费部署

1. 注册/登录 `app.koyeb.com`

2. 点击 `Create Service` -> `Docker` -> `Next`

3. `Image` 填写 `registry.gitlab.com/coiaprant/sleepy/dashboard:latest`, 点击 `Next`

4. 选择 `Free`, 地区 欧美 二选一, 点击 `Next`

---

5. 

找到 `Environment variables and files`

- `Name` 填写 `SECRET`
- `Value` 填写 你想设置的密钥

点击 `Add another`, 在新的一行

- `Name` 填写 `SITE_CONFIG`
- `Value` 填写 `site_config.json` 的 `base64` 格式

> `site_config.json` 参考内容

```json
{
    "title": "Sleepy",
    "description": "",
    "nickname": "TA"
}
```

> 推荐使用 [base64.us](https://base64.us)

找到 `Exposed ports`

- `Port` 填写 `80`
- `Protocol` 改为 `HTTP/2`

其他保持不变, 点击页面底部的 `Deploy`