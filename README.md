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


## 开发流程

git常用命令
```
# 检出最新的 develop 分支
$ git checkout develop
$ git pull

# 本地检出新的 feature 分支
$ git checkout -b feature/20250112_new_feature

# 可 Push 到远端备份 feature 分支的开发进度
$ git push --set-upstream origin feature/20250112_new_feature

# 设定 --set-upstream 后，后续 Push 动作直接在本分支使用 git push 命令即可完成
$ git push

# 拉取远程分支到本地，并切换到新分支
$ git checkout -b 本地分支名x origin/远程分支名x

# 拉取远程分支到本地

git fetch 本地分支名x origin/远程分支名x

git pull origin <远程分支名>:<本地分支名>
```


**一个功能对于一个feature分支，每个分支开发完后，需要先提PR，合并到testing分支。
严格禁止，直接merge到master分支。**