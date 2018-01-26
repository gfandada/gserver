package autoconversion

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	fmsgs, finner *os.File
)

var (
	handlerstr = "\n" +
		`func NewHandlers() {` + "\n" +
		"	// if you want to do with CLOSE_CONNECT, you have to achieve this function.\n" +
		"	// Services.Register(uint16(Services.CLOSE_CONNECT), CloseHandler)\n"
)

func Conversion(src, dest string) {
	load(src, dest)
}

func load(src, dest string) {
	os.RemoveAll(dest)
	os.Mkdir(dest, 777)
	dir_list, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}
	fmsgs, err = os.Create(dest + "/" + "msgs.go")
	if err != nil {
		return
	}
	defer fmsgs.Close()
	writeString :=
		"package out\n\n" +
			`import (
	Msg "protomsg"
	"github.com/gfandada/gserver/network"
)
` + "\n" +
			`func NewMsgCoder() *network.MsgManager {` + "\n" +
			"	coder := network.NewMsgManager()\n"
	io.WriteString(fmsgs, writeString)
	finner, err = os.Create(dest + "/" + "inner.go")
	if err != nil {
		return
	}
	defer finner.Close()
	writeString = "package out\n\n" +
		`import (
	Msg "protomsg"
	Services "github.com/gfandada/gserver/services"
	Session "github.com/gfandada/gserver/services/service"
	
	"github.com/gfandada/gserver/network"
)
` + "\n"
	io.WriteString(finner, writeString)
	for _, v := range dir_list {
		parse(src+v.Name(), dest)
	}
	writeString = "	return coder\n}\n"
	io.WriteString(fmsgs, writeString)
	io.WriteString(finner, handlerstr)
	writeString = "	return\n}\n"
	io.WriteString(finner, writeString)
}

func parse(file, dest string) error {
	fi, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	m := make(map[string][]map[string][]interface{})
	err = yaml.Unmarshal(fi, &m)
	if err != nil {
		return err
	}
	for key := range m {
		createHandler(dest, key, m[key])
	}
	return nil
}

func createHandler(dest, module string, args []map[string][]interface{}) error {
	f, err := os.Create(dest + "/" + module + "handler.go")
	if err != nil {
		return err
	}
	defer f.Close()
	writeString :=
		"package out\n\n" +
			`import (
	Msg "protomsg"
	Session "github.com/gfandada/gserver/services/service"
	
	//"github.com/gfandada/gserver/network"
	//"github.com/golang/protobuf/proto"
)
` + "\n"
	_, err = io.WriteString(f, writeString)
	if err != nil {
		return err
	}
	for _, v := range args {
		for key := range v {
			req := v[key][1].(string)
			reqid := v[key][0].(int)
			writeString =
				"" +
					"func " + req + "Handler(req *Msg." + req + ", sess *Session.Session" + ") interface{} {\n" +
					`	return nil
}` + "\n\n"
			io.WriteString(f, writeString)
			if len(v[key]) == 2 {
				writeString =
					"func Inner" + req + "Handler(args []interface{}" + ") []interface{} {\n" +
						"return []interface{}{" + req + "Handler(args[0].(*network.RawMessage).MsgData.(*Msg." + req +
						"), args[1].(*Session.Session))}" + "\n}\n"
			} else if len(v[key]) == 4 {
				ackid := v[key][2].(int)
				writeString =
					"func Inner" + req + "Handler(args []interface{}" + ") []interface{} {" +
						"ret:=" + req + "Handler(args[0].(*network.RawMessage).MsgData.(*Msg." + req +
						"), args[1].(*Session.Session))\n" +
						"return []interface{}{network.RawMessage{MsgId: uint16(" +
						fmt.Sprintf("%d", ackid) + "),MsgData:ret}}\n}\n"
			}
			handlerstr += "	Services.Register(uint16(" + fmt.Sprintf("%d", reqid) + "), Inner" + req + "Handler)\n"
			io.WriteString(finner, writeString)
			createMsgs(v[key])
			//createHandlers(key, v[key])
		}
	}
	return nil
}

func createMsgs(args []interface{}) error {
	writeString := ""
	if len(args) == 2 {
		reqid := args[0].(int)
		req := args[1].(string)
		writeString =
			"	coder.Register(&network.RawMessage{MsgId: uint16(" + fmt.Sprintf("%d", reqid) + "), MsgData: &Msg." + req + "{}})\n"
		io.WriteString(fmsgs, writeString)
	} else if len(args) == 4 {
		reqid := args[0].(int)
		req := args[1].(string)
		ackid := args[2].(int)
		ack := args[3].(string)
		writeString =
			"	coder.Register(&network.RawMessage{MsgId: uint16(" + fmt.Sprintf("%d", reqid) + "), MsgData: &Msg." + req + "{}})\n" +
				"	coder.Register(&network.RawMessage{MsgId: uint16(" + fmt.Sprintf("%d", ackid) + "), MsgData: &Msg." + ack + "{}})\n"
		io.WriteString(fmsgs, writeString)
	}
	return nil
}

//func createHandlers(key string, args []interface{}) error {
//	writeString := ""
//	if len(args) == 2 || len(args) == 4 {
//		reqid := args[0].(int)
//		req := args[1].(string)
//		writeString =
//			"	Services.Register(uint16(" + fmt.Sprintf("%d", reqid) + "), Inner" + req + "Handler)\n"
//		io.WriteString(fhandlers, writeString)
//	}
//	return nil
//}
