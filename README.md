# zhku-oj-server
ZhkuOJ项目服务端代码仓库

## 项目简介
```
.
├── README.md
├── air
├── cmd
│   ├── main.go
│   └── runServer.go
├── conf
│   ├── air.conf
│   ├── config.go
│   └── config.yaml
├── deploy
├── go.mod
├── go.sum
├── log
│   └── server.log
└── pkg
    ├── app
    │   └── api-server
    │       ├── common
    │       │   └── interface.go
    │       ├── dto
    │       │   └── dto_user.go
    │       ├── server
    │       │   ├── server.go
    │       │   └── server_user.go
    │       └── service
    │           ├── service.go
    │           └── service_user.go
    ├── dao
    │   ├── constant.go
    │   ├── dao.go
    │   └── dao_user.go
    ├── models
    │   ├── mongo.go
    │   └── mongo_user.go
    └── utils
        ├── common.go
        ├── config.go
        ├── constant.go
        ├── errors.go
        ├── file
        │   ├── file.go
        │   └── shutil.go
        ├── logger.go
        ├── middleware
        │   ├── auth.go
        │   ├── cors.go
        │   └── logger.go
        ├── mongo.go
        ├── redis.go
        └── response.go

```

## 项目启动

修改配置文件conf/config.yaml，把服务端口，地址，mongo地址修改成自己本地环境的地址

使用air进行热重启。
- 下载对于平台的air安装包，参考[air github](https://github.com/air-verse/air)，在release里面找到适合自己平台的压缩包
- 切换至当前路径，执行
  `./air -c conf/air.conf`


