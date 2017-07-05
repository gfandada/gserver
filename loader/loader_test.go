package loader

import (
	"testing"
)

type RedisCfg struct {
	Ip   string
	Port int
}

type TestCfg struct {
	Name    string
	Enabled bool
	Qq      int
	Numbers []int
	Email   []string
	Redis   RedisCfg
}

func Test_run(t *testing.T) {
	data := new(TestCfg)
	Loader("./test.json", data)
	if data.Email[1] != "gfandada@gmail.com" {
		t.Error("loader error")
	}
}
