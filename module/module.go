// 封装了模块拔插式的操作
package module

import (
	"fmt"
	"lib/logger"
	"os"
	"os/signal"
	"sync"
)

// 每个模块的详情
type module struct {
	Moduler   Imodule        // 模块实现的接口
	ChanClose chan bool      // 模块是否关闭的通道
	WaitSync  sync.WaitGroup // 保证子携程能执行完
}

// 使用切片作为模块的容器
var mods []*module

// 运行模块
func Run(mods ...Imodule) {
	for _, value := range mods {
		register(value)
	}
	initMods()
	chanSig := make(chan os.Signal, 1)
	signal.Notify(chanSig, os.Interrupt, os.Kill)
	sig := <-chanSig
	logger.Error(fmt.Sprintf("recv kill %v", sig))
	destroy()
}

// 注册模块
func register(iMod Imodule) {
	mod := new(module)
	mod.Moduler = iMod
	mod.ChanClose = make(chan bool, 1)
	mods = append(mods, mod)
}

// 回收模块
func destroy() {
	for i := len(mods) - 1; i >= 0; i-- {
		mod := mods[i]
		mod.ChanClose <- true
		mod.WaitSync.Wait()
		defer func() {
			if r := recover(); r != nil {
				logger.Error(fmt.Sprintf("destroy module panic error，%v", r))
			}
		}()
		mod.Moduler.OnDestroy()
	}
}

// 初始化
func initMods() {
	for _, value := range mods {
		value.Moduler.OnInit()
	}
	for _, value := range mods {
		value.WaitSync.Add(1)
		go run(value)
	}
}

// 运行模块
// FIXME 这里可以看出，插件模块的开发是有规范的，插件模块除了需要实现具体的接口以外
//       还需要满足一主多从的携程模式，主携程应该及时返回
func run(mod *module) {
	mod.Moduler.Run(mod.ChanClose)
	mod.WaitSync.Done()
}
