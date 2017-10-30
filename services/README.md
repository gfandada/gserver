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
优点：
1 对于游戏开发：特别适用于小服（小地图多副本）开发，对大服（无缝大世界）暂时支持不够，需要持续开发。
2 对于其他实时流业务，本框架都能很好的提供支持，也是gserver一开始的定位：Gen-Server-Engine
```
