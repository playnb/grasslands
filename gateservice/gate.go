package gateservice

import msg "github.com/playnb/mustang/msg.pb"

import (
	"github.com/golang/protobuf/proto"
	"github.com/playnb/mustang/gate"
	"github.com/playnb/mustang/log"
	"github.com/playnb/mustang/network"
	"github.com/playnb/mustang/network/protobuf"
	"github.com/playnb/mustang/utils"
	"strconv"
	"time"
)

//GateService实例
var Instance = newGateService()

func newGateService() *gateService {
	instance := new(gateService)
	instance.processor = protobuf.NewProtobufProcessor()
	instance.gateway = new(gate.WSGate)
	return instance
}

type gateService struct {
	processor   network.IProtobufProcessor
	gateway     *gate.WSGate
	SuperClient *msg.SuperServiceClient
}

//注册消息处理函数
func (g *gateService) RegisterHandle(msg proto.Message, msgHandler network.MsgHandler) {
	g.processor.Register(msg, msgHandler)
}

func (g *gateService) StartService(superRpcAddr string) {
	superClient, err := msg.DialSuperService("tcp", superRpcAddr)
	maxRetry := 10
	for err != nil {
		if maxRetry > 0 {
			maxRetry = maxRetry - 1
		} else {
			log.Fatal("连接SuperService失败")
			return
		}
		log.Error("连接SuperService失败,1秒后重试 :%v", err)
		time.Sleep(time.Second * 1)
		superClient, err = msg.DialSuperService("tcp", superRpcAddr)
	}
	res, err := superClient.Login(&msg.LoginRequst{ServiceIp: proto.String("127.0.0.1")})
	if err != nil {
		log.Fatal("[GATE] 登录SuperService失败 rpc:%s", superRpcAddr)
		return
	}
	g.SuperClient = superClient
	g.gateway.Addr = string(res.GetServiceIp()) + ":" + strconv.Itoa(int(res.GetExterPort()))
	g.gateway.HTTPTimeout = 3 * 60
	g.gateway.MaxConnNum = 1000
	g.gateway.PendingWriteNum = 1000
	g.gateway.ProtobufProcessor = g.processor
	log.Trace("[GATE] 网关服务在%s:%d 启动", string(res.GetServiceIp()), res.GetExterPort())
	g.gateway.Run(utils.CloseSig)
}

/*
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
	utils.ProtobufProcess.Register(new(Msg.EchoMsg), gateHandleMsg)

	gateway := new(gate.WSGate)
	gateway.Addr = "localhost:" + "3000"
	gateway.HTTPTimeout = 3 * 60
	gateway.MaxConnNum = 1000
	gateway.PendingWriteNum = 1000
	gateway.ProtobufProcessor = utils.ProtobufProcess
	gateway.Run(utils.CloseSig)
}
*/
