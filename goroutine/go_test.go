package goroutine

import (
	"fmt"
	"testing"
)

func Test_go(t *testing.T) {
	pid, err := Start(new(TestGo), true)
	if err != nil {
		t.Error(err)
	}
	//Stop(pid)
	Cast(pid, []interface{}{0001, "gfandada1", "coding"})
	Cast(pid, []interface{}{0002, "gfandada2", "coding"})
	fmt.Println(Count())
	Stop(pid)
	fmt.Println(Count())
	Cast(pid, []interface{}{0003, "gfandada3", "coding"})
	//Stop(pid)
	fmt.Println(Count())
	Start(new(TestGo), true)
	Start(new(TestGo), true)
	Start(new(TestGo), true)
	fmt.Println(Count())
	pid, err = Start(new(TestGo), false)
	if err != nil {
		t.Error(err)
	}
	ret, err1 := Call(pid, []interface{}{0004, "gfandada4", "coding"}, 1)
	if err == nil {
		fmt.Println("同步调用", ret)
	} else {
		t.Error(err1)
	}
}

type TestGo struct {
}

func (*TestGo) initGo() {
	fmt.Println("【init】welcome to gs, you can contact gfandada@gmail.com")
}

func (*TestGo) closeGo() {
	fmt.Println("【close】welcome to gs, you can contact gfandada@gmail.com")
}

func (*TestGo) handler(input []interface{}) []interface{} {
	id := input[0].(int)
	name := input[1].(string)
	job := input[2].(string)
	fmt.Printf("【handler】用户%s(%d)从事%s工作\n", name, id, job)
	return []interface{}{"gfandada@gmail.com"}
}
