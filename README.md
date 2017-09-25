# gserver
```
This project aims to provide a solution for real-time message flow. 
You can create GameServer or others, with gserver.
The communication protocol of GateWay has supported Tcp and Websocket.
DEMO is here : https://github.com/gfandada/gserver_demo.
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
4.add tcp gateway -- DONE
```
### Message
```
client->gateway
	----------------------------
	| len | seq | id | message |
	----------------------------
	len:seq + id + message，占用2个字节(uint16)
	seq:从1自增的序列号，占用4个字节(uint32)
	id:协议号，占用2个字节(uint16)
	message:业务数据，占用len-6字节，可以使用任意编码：pb/json等，本框架内置了pb编码器

gateway->client
	----------------------
	| len | id | message |
	----------------------
	len:id + message的长度，占用2个字节(uint16)
	id:协议号，占用两个字节(uint16)
	message:业务数据，占用len-2字节，可以使用任意编码：pb/json等，本框架内置了pb编码器
	
gateway<->service(base pb3)
	type Data_Frame struct {
		Type    Data_FrameType
		Message []byte
	}
```
### dataflow
![image](https://github.com/gfandada/gserver/blob/master/png/dataflow.png)
