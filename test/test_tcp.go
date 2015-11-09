package test

import (
	"TestGo/Msg"
	"github.com/playnb/mustang/log"
	"github.com/playnb/mustang/network"
	"github.com/playnb/mustang/network/protobuf"
	"github.com/playnb/mustang/utils"
	"github.com/golang/protobuf/proto"
	"reflect"
	"time"
)

///////////////////////////////////////////////////////////////////////////////

func tcpHandleMsg(agent network.IAgent, msgType reflect.Type, message interface{}, data []interface{}) {
	if "Msg.EchoMsg" != msgType.String() {
		return
	}
	msg := message.(*Msg.EchoMsg)
	log.Debug(reflect.TypeOf(agent).String() + ":" + msg.GetEchoString())
	//msg.EchoString = proto.String("i will be back.....")
	//agent.WriteMsg(msg)
}

///////////////////////////////////////////////////////////////////////////////
type TcpServer struct {
	Addr              string
	agent             *TCPServerAgent
	ProtobufProcessor *protobuf.Processor
}

func (s *TcpServer) Run(closeSig chan bool) {
	server := &network.TCPServer{}
	server.Addr = s.Addr
	server.MaxConnNum = 1000
	server.PendingWriteNum = 1000
	server.NewAgent = func(conn *network.TCPConn) network.IAgent {
		a := new(TCPServerAgent)
		a.server = s
		a.SetConn(conn)
		a.SetName("TCPServerAgent")
		a.SetProtobufProcessor(utils.ProtobufProcess)
		s.agent = a
		return a
	}
	log.Trace("TCP服务器启动")
	server.Start()
	//等待close信号
	<-closeSig
	server.Close()
}

func (s *TcpServer) Close() {
}

type TCPServerAgent struct {
	network.TCPAgent
	server *TcpServer
}

///////////////////////////////////////////////////////////////////////////////

const testAddr = "127.0.0.1:987"

type TCPClientAgent struct {
	network.TCPAgent
}

func runClient() {
	client := new(network.TCPClient)
	client.Addr = testAddr
	client.MaxTry = 3
	client.ConnectInterval = 3 * time.Second
	client.NewAgent = func(conn *network.TCPConn) network.IAgent {
		a := new(TCPClientAgent)
		a.SetConn(conn)
		a.SetProtobufProcessor(utils.ProtobufProcess)
		a.SetName("TCPClientAgent")
		return a
	}
	client.Start()

	msg := new(Msg.EchoMsg)
	msg.EchoString = proto.String("我发消息了")

	for {
		time.Sleep(3 * time.Second)
		client.WriteMsg(msg)
	}
}

func TestTcp() {
	utils.ProtobufProcess.Register(new(Msg.EchoMsg), tcpHandleMsg)

	go runClient()

	time.Sleep(10 * time.Second)
	server := new(TcpServer)
	server.ProtobufProcessor = utils.ProtobufProcess
	server.Addr = testAddr
	server.Run(utils.CloseSig)
}
