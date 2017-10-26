package loader

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gfandada/gserver/logger"
)

func Test_loader(t *testing.T) {
	logger.Start("../gservices/test.xml")
	Init("./test/")
	loader := new(Loader)
	if ret, _ := loader.Get("Equipment", 3, "Price"); ret != uint32(500) {
		fmt.Println(reflect.TypeOf(ret))
		t.Error("get error")
	}
	if ret, _ := loader.Get("Equipment", 7, "SalePrice"); ret != uint32(1000) {
		t.Error("get error")
	}
	if ret, _ := loader.Get("Equipment", 8, "Icon"); ret != "e_icon" {
		t.Error("get error")
	}
	// 获取关联数据
	//	ret, _ := loader.GetCorrelation("buildinginfo", 1003, "userlevel")
}
