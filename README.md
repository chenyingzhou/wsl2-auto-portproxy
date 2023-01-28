# wsl2-tcpproxy
用于将Windows的TCP端口转发至wsl2，另外，也支持通过配置转发至wsl2以外的端口

## 转发至WSL2
wsl2需要安装`net-tools`，Ubuntu可通过以下命令安装
```bash
sudo apt-get install net-tools
```

## 如何安装
前往 [release](https://github.com/chenyingzhou/wsl2-tcpproxy/releases) 下载最新版本的`wsl2-tcpproxy.exe`
#### 或从源码构建
```bash
make build
```
#### 如果环境为Windows，用以下命令替代`make build`
```bash
go build -o ./wsl2-tcpproxy.exe
```

## 配置文件
文件位于`$HOME/.wsl2-tcpproxy.json`，首次运行时会自动创建   
示例：
```json
{
  "predefined": [
    "6380:6379"
  ],
  "ignore": [
    "443"
  ],
  "custom": [
    "8081:192.168.1.99:8080"
  ]
}
```
- predefined: 预设的端口映射，格式为`winport:wslport`，该设置仅针对`winport和wslport不一致`的情况，端口相同的无需设置
- ignore: 忽略的wsl2的端口，该端口不会被代理
- custom: wsl2以外的代理，支持转发到任意机器的任意端口，格式为`winport:remoteip:remoteport`

## License
MIT

