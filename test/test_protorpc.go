package test

import (
	"TestGo/Msg"
	"github.com/playnb/mustang/log"
	//"github.com/chai2010/protorpc"
	"github.com/golang/protobuf/proto"
	"strconv"
	"sync"
	"time"
)

type EchoService struct {
	sync.Mutex
	count uint32
}

func (t *EchoService) Echo(args *Msg.EchoRequest, reply *Msg.EchoResponse) error {
	//t.Lock()
	//defer t.Unlock()
	reply.Msg = proto.String("Echo:" + args.GetMsg())
	t.count++
	log.Trace(strconv.Itoa(int(t.count)))
	time.Sleep(time.Second * 4)
	return nil
}

func (t *EchoService) EchoTwice(args *Msg.EchoRequest, reply *Msg.EchoResponse) error {
	reply.Msg = proto.String("EchoTwice:" + args.GetMsg() + args.GetMsg())
	return nil
}

func startProtorpcService() {
	Msg.ListenAndServeEchoService("tcp", NetworkAddr, new(EchoService))
}

func TestProtorpc() {
	go startProtorpcService()

	echoClient, err := Msg.DialEchoService("tcp", NetworkAddr)
	for err != nil {
		log.Error("连接服务失败 :%v", err)
		time.Sleep(time.Second * 1)
		echoClient, err = Msg.DialEchoService("tcp", NetworkAddr)
	}
	echoClient, err = Msg.DialEchoService("tcp", NetworkAddr)

	wg := new(sync.WaitGroup)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Debug("RPC Call Begin..." + strconv.Itoa(i))
			reply, err := echoClient.Echo(&Msg.EchoRequest{Msg: proto.String("Hello!" + strconv.Itoa(i))})
			if err != nil {
				log.Error("EchoTwice: %v", err)
			}
			log.Trace(reply.GetMsg())
			log.Debug("RPC Call Finish..." + strconv.Itoa(i))
		}(i)
	}
	/*
		client, err := protorpc.Dial("tcp", `127.0.0.1:9527`)
		if err != nil {
			log.Fatal("protorpc.Dial: %v", err)
		}
		defer client.Close()

		echoClient1 := &Msg.EchoServiceClient{client}
		echoClient2 := &Msg.EchoServiceClient{client}
		reply, err = echoClient1.Echo(args)
		log.Trace(reply.GetMsg())
		reply, err = echoClient2.EchoTwice(args)
		log.Trace(reply.GetMsg())
		_, _ = reply, err
	*/

	wg.Wait()
	echoClient.Close()
}
