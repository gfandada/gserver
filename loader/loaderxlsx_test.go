package loader

import (
	"fmt"
	"testing"

	"github.com/gfandada/gserver/logger"
)

func Test_loader(t *testing.T) {
	logger.Start("../gservices/test.xml")
	Init("./test/")
	loader := new(Loader)
	if ret, _ := loader.Get("buildinginfo", 1003, "posy"); ret != -2002 {
		t.Error("get error")
	}
	if ret, _ := loader.Get("buildinginfo", 1003, "posx"); ret != 1001 {
		t.Error("get error")
	}
	if ret, _ := loader.Get("buildinginfo", 1003, "info"); ret != "嗯嗯嗯嗯呃呃123" {
		t.Error("get error")
	}
	if ret, _ := loader.Get("userlevel", 4, "info"); ret != "嗯嗯嗯嗯呃呃123" {
		t.Error("get error")
	}
	// 获取关联数据
	//	ret, _ := loader.GetCorrelation("buildinginfo", 1003, "userlevel")
}
