package test

import (
	"TestGo/Msg"
	"github.com/golang/protobuf/proto"
	"github.com/playnb/mustang/gate"
	"github.com/playnb/mustang/log"
	"github.com/playnb/mustang/network"
	"github.com/playnb/mustang/network/protobuf"
	"github.com/playnb/mustang/utils"
	"reflect"
)

type GateService struct {
}

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

func TestGate() {
	protobuf.ProtobufProcess.Register(new(Msg.EchoMsg), gateHandleMsg)

	gateway := new(gate.WSGate)
	gateway.Addr = "localhost:" + "3000"
	gateway.HTTPTimeout = 3 * 60
	gateway.MaxConnNum = 1000
	gateway.PendingWriteNum = 1000
	gateway.ProtobufProcessor = protobuf.ProtobufProcess
	gateway.Run(utils.CloseSig)
}
