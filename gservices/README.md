### gservices包介绍
```
gservices包是以erlang:gen_server/gen_timer....为目标所做的封装。
genserver:通用的请求-应答式服务器模型。
gentimer:通用的定时器服务器模型。
不足与缺陷：
1 相对于erlang:gen_server/gen_timer功能略有不足，会持续优化。
2 暂时未实现link,restart,name等特性。
3 纯内存，没有持久化。
```