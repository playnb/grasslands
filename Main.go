package main

import (
	"TestGo/Msg"
	"github.com/playnb/mustang/auth"
	"github.com/playnb/mustang/gate"
	"github.com/playnb/mustang/log"
	"github.com/playnb/mustang/network"
	//"TestGo/mustang/network/protobuf"
	"github.com/playnb/mustang/utils"
	"TestGo/test"
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
)

type NeedIF interface {
	IF()
}

type EchoRequest struct {
	Msg              *string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (this *EchoRequest) IF() {

}

func WhoAmI(me interface{}) string {
	t := reflect.TypeOf(me)
	return t.String()[1:]
}

func test_WhoAmI() {
	echo := &EchoRequest{}
	log.Trace(WhoAmI(echo))
}

func pbyte(b []byte) {
	fmt.Print("OUT: ")
	for i := 0; i < len(b); i++ {
		fmt.Print(b[i], ",")
	}
	fmt.Print("\n")
}

func cbyte(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = b[i] + 1
	}
}

func test_slice() {
	b := [10]byte{}
	for i := 0; i < 10; i++ {
		b[i] = byte(i)
	}
	pbyte(b[0:])
	cbyte(b[0:])
	pbyte(b[0:])
}

/*
func handleMsg(t reflect.Type, m []interface{}) {
	if "Msg.EchoMsg" != t.String() {
		return
	}
	msg := m[0].(*Msg.EchoMsg)
	log.Debug(msg.GetEchoString())
}

func test_process() {
	utils.ProtobufProcess.Register(new(Msg.EchoMsg), handleMsg)

	msg1 := new(Msg.EchoMsg)
	msg1.EchoCode = proto.Uint32(1)
	msg1.EchoString = proto.String("hi all")
	if data, err := utils.ProtobufProcess.PackMsg(msg1); err != nil {
		log.Trace(err.Error())
	} else {
		utils.ProtobufProcess.Handler(data)
	}
}
*/

func gateHandleMsg(agant network.IAgent, msgType reflect.Type, message interface{}, data []interface{}) {
	if "Msg.EchoMsg" != msgType.String() {
		return
	}
	msg := message.(*Msg.EchoMsg)
	log.Debug(msg.GetEchoString())
	msg.EchoString = proto.String("i will be back.....")
	agant.WriteMsg(msg)
	msg.EchoString = proto.String("i am back.....")
	agant.WriteMsg(msg)
}

func test_gate() {
	utils.ProtobufProcess.Register(new(Msg.EchoMsg), gateHandleMsg)

	gateway := new(gate.WSGate)
	gateway.Addr = "localhost:" + "3000"
	gateway.HTTPTimeout = 3 * 60
	gateway.MaxConnNum = 1000
	gateway.PendingWriteNum = 1000
	gateway.ProtobufProcessor = utils.ProtobufProcess
	gateway.Run(utils.CloseSig)
}

type IA interface {
	getA() int
}

type A struct {
	value_A int
}

func (a *A) dumpA() {
	log.Trace("A%d", a.value_A)
}
func (a *A) getA() int {
	return a.value_A
}

type B struct {
	IA
	value_B int
}

func (b *B) dumpB() {
	log.Trace("A%d B%d", b.getA(), b.value_B)
}

func main() {

	auth.InitAuthHttpService()

	log.Trace("启动......")

	test.TestSinaAuth()

	log.Trace("结束......")

	<-utils.CloseSig
	return
	test.TestProtorpc()
	utils.TestSnow()
}
