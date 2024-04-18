## 前沿
> - 本人主要是从事PHP语言开发，因2022年疫情时期在上海封控2个月，在一次偶然的机会接触到Golang，也算是初学Go语言，变学边记，一直坚持到现在，学完了Go的基础、语法、常见的包、网络协议、goroutine等知识，想上手结合Gin框架写一个简单的IM通讯系统，融会贯通一下Go的知识点。

## 项目介绍
> IM(即时通讯)，支持单聊，群聊，技术栈涉及到Websocket协议，MongoDB数据库，Gin框架的操作等。

## 核心包
```
https://github.com/gorilla/websocket
https://github.com/mongodb/mongo-go-driver
```

## 扩展安装
```
go get -u githun.com/gin-gonic/gin
go get github.com/gorilla/websocket
go get go.mongodb.org/mongo-driver/mongo
go get github.com/dgrijalva/jwt-go
```

## Docker 安装 MongoDB
```
docker run -d --name="im-mongo" -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin -p 27017:27017 mongo
```
