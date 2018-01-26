# 自动转换
autoconversion包是为了简化开发者的工作，主要是将协议id，协议结构，handler的注册和映射脚本化统一处理，
不需要额外的用户代码植入，对项目管理来说有很好的帮助。

# 语法介绍
使用标准yaml文件管理
以sample.yml为例说明：  
login/test:【这个字段是自定义的模块名，会被自动生成一个loginhandler.go/testhandler.go的源文件】

1.标准的请求-应答模式【同步响应】  
LoginReq: [1,LoginReq,2,LoginAck]  
LoginReq【是handler函数名的前缀，可以自定义，会自动生成一个LoginReqHandler的端口函数】  
[1,LoginReq,2,LoginAck]【1表示client发出的协议id,LoginReq是请求的数据对象,2表示需要回复的协议id,LoginAck是需要回复的数据对象】  

2.标准的请求-无应答模式【无响应】  
TestReq: [3,TestReq]  
TestReq【是handler函数名的前缀，可以自定义，会自动生成一个TestReqHandler的端口函数】  
[3,TestReq]【1表示client发出的协议id,LoginReq是请求的数据对象,此请求不需要服务器响应】  

2.标准的请求-推送模式【异步响应】  
Test1Req: [4,Test1Req,5,Test1Push]  
类似于1  

# 注意事项
自动生成的代码目录需求适配自己的项目
