package test

import (
	"TestGo/Rpc/demo"
	"github.com/playnb/mustang/log"
	"git.apache.org/thrift.git/lib/go/thrift"
	"time"
)

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

type RpcServiceImpl struct {
}

func (this *RpcServiceImpl) Ping() error {
	log.Trace("有人ping我")
	return nil
}

func (this *RpcServiceImpl) Add(num1 int32, num2 int32) (int32, error) {
	log.Trace("有人add %v+%v", num1, num2)
	time.Sleep(1 * time.Second)
	return (num1 + num2), nil
}

func serverStub() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	//protocolFactory := thrift.NewTCompactProtocolFactory()

	serverTransport, err := thrift.NewTServerSocket(NetworkAddr)
	if err != nil {
		log.Error("Error!", err)
		return
	}

	handler := &RpcServiceImpl{}
	processor := demo.NewRpcServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	log.Trace("thrift server in", NetworkAddr)
	server.Serve()
}

func clientCall(min int32) {
	startTime := currentTimeMillis()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(NetworkAddr)
	for err != nil {
		transport, err = thrift.NewTSocket(NetworkAddr)
		if err != nil {
			log.Error("error resolving address:", err)
		}
		time.Sleep(1 * time.Second)
	}

	useTransport := transportFactory.GetTransport(transport)
	client := demo.NewRpcServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		log.Error("Error opening socket to 127.0.0.1:19090", " ", err)
		return
	}
	defer transport.Close()

	for i := min; i < min+3; i++ {
		r1, e1 := client.Add(i, i+1)
		log.Trace("%d %s %v %v", i, "Call->", r1, e1)
	}

	endTime := currentTimeMillis()
	log.Trace("Program exit. time->", endTime, startTime, (endTime - startTime))
}

func TestThrift() {
	go clientCall(100)
	go clientCall(200)
	go clientCall(300)
	serverStub()
}
