# gserver
```
gen tcp/websocket server base on message
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
```
### TODO
```
1.optimize handler exec loop
2.optimize conn pool
3.optimize inner panic
4.storage timer job base on gentimer
5.romote rpc for cluster
```
### cluster
![image](https://github.com/gfandada/gserver/blob/master/png/cluster.png)
### dataflow
![image](https://github.com/gfandada/gserver/blob/master/png/dataflow.png)
