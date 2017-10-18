### services包介绍
```
services包是微服务开发的基础组件。
gateway:支持通用的tcp/websocket。
service:支持通用的微服务。
discovery:服务注册和发现(基于本地配置no-watching)。
不足与缺陷：
1 服务发现、服务治理比较粗糙。
2 消息路由未采用较流行的类url，而是直接采用消息id分段机制。
3 由于gateway针对每条连接都是新增一个go进程作为目标service的路由进程，所以暂时只内置了一个
router进程。
```