# t2

## 项目简介

MCBOT 是一个用 Go 语言编写的项目，致力于帮助服主使用机器人进行管理服务器

## 目录结构

```
go.mod           # Go 依赖管理文件
go.sum           # Go 依赖校验文件
main.go          # 主程序入口
MCWhitelist.db   # 白名单数据库
README.MD        # 项目说明文档

event/           # 事件相关代码
  event.go
Message/         # 消息相关代码
  message.go
whitelist/       # 白名单相关代码
  whitelist.go
minecraft/       # TODO实现大部分minecraft服务器抽象
```

## 依赖管理

本项目使用 Go Modules 进行依赖管理。

安装依赖：
```bash
go mod tidy
```

## 构建与运行

编译项目：
```bash
go build -o t2.exe
```

运行项目：
```bash
go run main.go
```

## 贡献

欢迎提交 issue 和 pull request。

## 许可证

MIT License
