package goroutine

import (
	"fmt"
	"testing"
	"time"
)

func Test_go(t *testing.T) {
	pid, err := Start(new(Test))
	if err != nil {
		t.Error(err)
		return
	}
	go Cast(pid, "+", []interface{}{1, 2})
	add := func() {
		ret, err := Call(pid, "+", []interface{}{1, 2}, 2)
		if err != nil {
			fmt.Println(err.Error())
			t.Error(err)
			return
		}
		fmt.Println("1 + 2 = ", ret[0])
		if ret[0] != 3 {
			t.Errorf("1 + 2 != 3 ?????????????")
			return
		}
	}
	go add()
	sub := func() {
		ret, err := Call(pid, "-", []interface{}{1, 2}, 2)
		if err != nil {
			fmt.Println(err.Error())
			t.Error(err)
			return
		}
		fmt.Println("1 - 2 = ", ret[0])
		if ret[0] != -1 {
			t.Errorf("1 - 2 != -1 ?????????????")
			return
		}
	}
	go sub()
	mul := func() {
		ret, err := Call(pid, "*", []interface{}{1, 2}, 2)
		if err != nil {
			fmt.Println(err.Error())
			t.Error(err)
			return
		}
		fmt.Println("1 * 2 = ", ret[0])
		if ret[0] != 2 {
			t.Errorf("1 * 2 != 2 ?????????????")
			return
		}
	}
	go mul()
	time.Sleep(1e9)
	StopById(pid)
}

/*****************************实现进程装载器********************************/

type Test struct {
}

func (t *Test) name() string {
	return ""
}

func (t *Test) initGo() {
	fmt.Println("init..............")
}

func (t *Test) handler(msg string, args []interface{}, ret chan []interface{}) {
	fmt.Println("handler..............")
	// 异步的嘛
	if ret == nil {
		//...........do something...........
		return
	}
	// 同步的嘛
	switch msg {
	case "+":
		ret <- []interface{}{args[0].(int) + args[1].(int)}
	case "-":
		ret <- []interface{}{args[0].(int) - args[1].(int)}
	case "*":
		ret <- []interface{}{args[0].(int) * args[1].(int)}
	}
}

func (t *Test) closeGo() {
	fmt.Println("close..............")
}
