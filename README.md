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
current version is v0.1.3
next version will focus on:
1.optimize cluster data flow
2.optimize inner panic
3.storage timer job base on gentimer
```
### cluster
```
TODO MQ
```
![image](https://github.com/gfandada/gserver/blob/master/png/cluster.png)
### dataflow
```
TODO go exec
```
![image](https://github.com/gfandada/gserver/blob/master/png/dataflow.png)
