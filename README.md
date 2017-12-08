# gserver
This project aims to provide a solution for real-time message flow.  
You can create GameServer or others, with gserver.  
The communication protocol of GateWay has supported Tcp and Websocket.  
[DEMO](https://github.com/gfandada/gserver_demo "gserver-demo") is here. (tcp-gateway and mobagame-service)  
[LICENSE](LICENSE "Apache License 2.0") is here.
### Installation
golang version 1.9.2  
go get github.com/golang/protobuf/proto  
go get github.com/gorilla/websocket  
go get github.com/cihub/seelog  
go get github.com/koding/multiconfig  
go get github.com/garyburd/redigo/redis  
go get github.com/HuKeping/rbtree  
go get github.com/tealeg/xlsx  
go get google.golang.org/grpc  
go get github.com/go-sql-driver/mysql  
### TODO
current version is v0.8.3  
next version-v0.9.0 will focus on:  
1.optimize safe -- ING   
2.optimize microservice   
3.add inner logger -- ING   
4.add tcp gateway -- DONE   
5.add game util package(aoi,space,entity....) -- ING   
### CONF
```
更全的配置请查看demo工程
{
	"MaxConnNum": 2048, // 最大连接数:多余的连接将不会响应
	"PendingNum": 100,  // gateway->client异步ipc队列上限
	"MaxMsgLen": 1024,  // client<->gateway message上限:单位byte
	"MinMsgLen": 0,     // client<->gateway message下限:单位byte
	"ReadDeadline":60,  // gateway->client读超时:单位s
	"WriteDeadline":60, // gateway->client写超时:单位s
	"ServerAddress": "localhost:9527", // gateway地址
	"MaxHeader":1024,   // header上限(for websocket):单位byte
	"HttpTimeout": 10,  // http-get超时(for websocket):单位s
	"CertFile": "",     // for ssl
	"KeyFile": "",      // for ssl
	"Rpm":100,          // client->gateway流量上限:每分钟收到的报文数量上限
	"AsyncMQ":64,       // service->gateway异步ipc队列上限
	"GateWayIds":1999   // gateway本地路由id段(当前路由规则是简单的id分段规则)
}
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
	message:业务数据，占用len-6字节，可以使用任意编码：pb/json等，本框架内置了pb2编码器

gateway->client
	----------------------
	| len | id | message |
	----------------------
	len:id + message的长度，占用2个字节(uint16)
	id:协议号，占用两个字节(uint16)
	message:业务数据，占用len-2字节，可以使用任意编码：pb/json等，本框架内置了pb2编码器
	
gateway<->service(base pb3)
	type Data_Frame struct {
		Type    Data_FrameType
		Message []byte
	}
```
### dataflow
![image](https://github.com/gfandada/gserver/blob/master/png/dataflow.png)
