# Windows 监控端使用教程

前往 [GitLab Releases](https://gitlab.com/CoiaPrant/Sleepy/-/releases)

找到对应设备的 `sleepy-agent` 下载并解压

随便找个地方新建文件夹

将 `sleepy-agent.exe` 放到文件夹内

在文件夹内创建 `sleepy.bat`

填写以下内容 并替换 `api` 和 `secret`

```shell
@echo off
sleepy-agent.exe --api example.com:443 --secret Example
```

然后双击 `sleepy.bat` 启动程序

> 目前没找到自启动的方法

## 参数解析

### 总览

```shell
  -api string
        RPC api address
  -debug
        Show debug logs
  -device string
        Report device name (default: follow-system)
  -h    Show help
  -insecure
        Disable SSL/TLS certificate verification
  -interval int
        Report state interval (default 2)
  -log string
        Log file location (default: stdout)
  -rpc string
        RPC type (valid: google) (default "google")
  -secret string
        RPC secret
  -v    Show version
```

### 说明

- `api`

面板地址, 例如 `--api example.com:443`

> 请不要带 http/https, 务必按例子格式填写

- `secret`

密钥, 需和面板端设置一致, 例如 `--secret Example`

- `device`

设备名称, 留空获取系统的主机名称, 例如 `--device "My Computer"`

- `interval`

上报周期, 默认2秒, 例如 `--interval 2`