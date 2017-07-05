// 封装了装载器，统一加载配置文件
package loader

import (
	"github.com/koding/multiconfig"
)

// 配置装载器
func Loader(path string, data interface{}) {
	multiconfig.NewWithPath(path).MustLoad(data)
}
