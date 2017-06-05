// 封装了集群layer的操作
package cluster

var GatewayLayer string = "gateway" // 网关层标识
var CalcLayer string = "calc"       // 计算层标识
var DbLayer string = "db"           // 数据层标识
var LogLayer string = "log"         // 日志层标识

var Rand string = "rand"   // 随机计算节点权重，在权重高的节点上执行handle
var Load string = "load"   // 根据运行时负载计算节点权重，在权重高的节点上执行handle
var Local string = "local" // 在本节点上执行handle

// 随机权重算法
func RandNode() {
}

// 负载权重算法
func LoadNode() {
}

// 本地权重算法
func LocalNode() {
}
