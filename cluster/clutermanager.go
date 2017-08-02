package cluster

import (
	"sync"

	"github.com/gfandada/gserver/loader"
	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/util"
	"google.golang.org/grpc"
)

type client struct {
	key    string
	conn   *grpc.ClientConn
	weight uint32 // 预留权重值
}

// 服务结构
type service struct {
	clients []client
	weight  uint32 // 预留权重值
}

// 服务池
type service_pool struct {
	names    map[string]bool     // 服务名：主要用于检查
	services map[string]*service // 服务详情
	mu       sync.RWMutex
}

var (
	_default_pool service_pool
	once          sync.Once
)

type instance struct {
	Id      string
	Address string
}

type cservice struct {
	Name     string
	Instance []instance
}

type cluster struct {
	Services []cservice
}

// 初始化
func Init(path string) {
	once.Do(func() { _default_pool.init(path) })
}

func (p *service_pool) init(path string) {
	p.names = make(map[string]bool)
	p.services = make(map[string]*service)
	s := new(cluster)
	loader.LoadJson(path, s)
	for _, v := range s.Services {
		for _, i := range v.Instance {
			p.addService(v.Name, i.Id, i.Address)
		}
	}
}

//  添加服务
func (p *service_pool) addService(serviceName, key, value string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.names[serviceName] = true
	instance := &service{}
	if conn, err := grpc.Dial(value, grpc.WithBlock(), grpc.WithInsecure()); err == nil {
		instance.clients = append(instance.clients, client{key, conn, 0})
		logger.Info("cluster service added succeed: %s -> %s", key, value)
		// TODO 回调通知
	} else {
		logger.Info("cluster service added error: %s -> %s", key, value)
	}
	p.services[serviceName] = instance
}

// 移除服务
func (p *service_pool) removeService(serviceName, key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.names[serviceName] {
		return
	}
	service := p.services[serviceName]
	if service == nil {
		logger.Error("no such service: %s", serviceName)
		return
	}
	for k := range service.clients {
		if service.clients[k].key == key {
			service.clients[k].conn.Close()
			service.clients = append(service.clients[:k], service.clients[k+1:]...)
			logger.Info("service removed: %s-%s", serviceName, key)
			return
		}
	}
}

// 通过名字获取一个服务的实例
func (p *service_pool) get_service(serviceName string) (conn *grpc.ClientConn) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if _, ok := p.names[serviceName]; !ok {
		return nil
	}
	service := p.services[serviceName]
	if service == nil {
		return nil
	}
	leng := len(service.clients)
	if leng == 0 {
		return nil
	}
	// TODO 先随机吧
	id := int(util.RandInterval(int32(0), int32(leng-1)))
	return service.clients[id].conn
}

// 通过获取所有服务的实例
func (p *service_pool) get_services() (conn map[string]*grpc.ClientConn) {
	for serviceName := range p.names {
		conn[serviceName] = p.get_service(serviceName)
	}
	return
}

// 获取一个服务
func GetService(service string) *grpc.ClientConn {
	return _default_pool.get_service(service)
}

// 获取所有服务
func GetServices() map[string]*grpc.ClientConn {
	return _default_pool.get_services()
}

// 获取所有服务名
func GetServiceNames() map[string]bool {
	return _default_pool.names
}
