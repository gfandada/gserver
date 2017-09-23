# gserver
```
This project aims to provide a solution for real-time message flow. You can create GameServer or others, with gserver.
The communication protocol is based on Websocket.
```
### Installation
```
go get github.com/golang/protobuf/proto
go get github.com/gorilla/websocket
go get github.com/cihub/seelog
go get github.com/koding/multiconfig
go get github.com/garyburd/redigo/redis
go get github.com/HuKeping/rbtree
go get github.com/tealeg/xlsx
go get google.golang.org/grpc
go get github.com/go-sql-driver/mysql
```
### TODO
```
current version is v0.8.3
next version-v0.9.0 will focus on:
1.optimize safe
2.optimize microservice
3.add inner logger
4.add tcp gateway
```
### dataflow
![image](https://github.com/gfandada/gserver/blob/master/png/dataflow.png)
