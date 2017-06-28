package loader

import (
	"testing"
)

func Test_run(t *testing.T) {
	ch := make(chan bool, 1)
	ch <- true
	loader := new(Loader)
	loader.OnInit()
	loader.Run(ch)
}
