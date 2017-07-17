// 加载json配置文件
package loader

import (
	"github.com/koding/multiconfig"
)

// 配置装载器
func LoadJson(path string, data interface{}) {
	multiconfig.NewWithPath(path).MustLoad(data)
}
