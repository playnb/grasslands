// Code generated by protoc-gen-go.
// source: Msg/message.proto
// DO NOT EDIT!

/*
Package Msg is a generated protocol buffer package.

It is generated from these files:
	Msg/message.proto

It has these top-level messages:
	EchoMsg
	EchoMsg2
	EchoRequest
	EchoResponse
*/
package Msg

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import "io"
import "log"
import "net"
import "net/rpc"
import "time"
import protorpc "github.com/chai2010/protorpc"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type EchoMsg struct {
	EchoCode         *uint32  `protobuf:"varint,1,opt,name=echo_code,def=77" json:"echo_code,omitempty"`
	EchoString       *string  `protobuf:"bytes,2,opt,name=echo_string" json:"echo_string,omitempty"`
	EchoVec          []string `protobuf:"bytes,3,rep,name=echo_vec" json:"echo_vec,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *EchoMsg) Reset()         { *m = EchoMsg{} }
func (m *EchoMsg) String() string { return proto.CompactTextString(m) }
func (*EchoMsg) ProtoMessage()    {}

const Default_EchoMsg_EchoCode uint32 = 77

func (m *EchoMsg) GetEchoCode() uint32 {
	if m != nil && m.EchoCode != nil {
		return *m.EchoCode
	}
	return Default_EchoMsg_EchoCode
}

func (m *EchoMsg) GetEchoString() string {
	if m != nil && m.EchoString != nil {
		return *m.EchoString
	}
	return ""
}

func (m *EchoMsg) GetEchoVec() []string {
	if m != nil {
		return m.EchoVec
	}
	return nil
}

type EchoMsg2 struct {
	EchoString       *string `protobuf:"bytes,1,opt,name=echo_string" json:"echo_string,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *EchoMsg2) Reset()         { *m = EchoMsg2{} }
func (m *EchoMsg2) String() string { return proto.CompactTextString(m) }
func (*EchoMsg2) ProtoMessage()    {}

func (m *EchoMsg2) GetEchoString() string {
	if m != nil && m.EchoString != nil {
		return *m.EchoString
	}
	return ""
}

// ///////////////////////////////////////////
type EchoRequest struct {
	Msg              *string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *EchoRequest) Reset()         { *m = EchoRequest{} }
func (m *EchoRequest) String() string { return proto.CompactTextString(m) }
func (*EchoRequest) ProtoMessage()    {}

func (m *EchoRequest) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

type EchoResponse struct {
	Msg              *string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *EchoResponse) Reset()         { *m = EchoResponse{} }
func (m *EchoResponse) String() string { return proto.CompactTextString(m) }
func (*EchoResponse) ProtoMessage()    {}

func (m *EchoResponse) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

type EchoService interface {
	Echo(in *EchoRequest, out *EchoResponse) error
	EchoTwice(in *EchoRequest, out *EchoResponse) error
}

// AcceptEchoServiceClient accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func AcceptEchoServiceClient(lis net.Listener, x EchoService) {
	srv := rpc.NewServer()
	if err := srv.RegisterName("EchoService", x); err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
		}
		go srv.ServeCodec(protorpc.NewServerCodec(conn))
	}
}

// RegisterEchoService publish the given EchoService implementation on the server.
func RegisterEchoService(srv *rpc.Server, x EchoService) error {
	if err := srv.RegisterName("EchoService", x); err != nil {
		return err
	}
	return nil
}

// NewEchoServiceServer returns a new EchoService Server.
func NewEchoServiceServer(x EchoService) *rpc.Server {
	srv := rpc.NewServer()
	if err := srv.RegisterName("EchoService", x); err != nil {
		log.Fatal(err)
	}
	return srv
}

// ListenAndServeEchoService listen announces on the local network address laddr
// and serves the given EchoService implementation.
func ListenAndServeEchoService(network, addr string, x EchoService) error {
	lis, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	srv := rpc.NewServer()
	if err := srv.RegisterName("EchoService", x); err != nil {
		return err
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
		}
		go srv.ServeCodec(protorpc.NewServerCodec(conn))
	}
}

type EchoServiceClient struct {
	*rpc.Client
}

// NewEchoServiceClient returns a EchoService stub to handle
// requests to the set of EchoService at the other end of the connection.
func NewEchoServiceClient(conn io.ReadWriteCloser) *EchoServiceClient {
	c := rpc.NewClientWithCodec(protorpc.NewClientCodec(conn))
	return &EchoServiceClient{c}
}

func (c *EchoServiceClient) Echo(in *EchoRequest) (out *EchoResponse, err error) {
	if in == nil {
		in = new(EchoRequest)
	}
	out = new(EchoResponse)
	if err = c.Call("EchoService.Echo", in, out); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *EchoServiceClient) EchoTwice(in *EchoRequest) (out *EchoResponse, err error) {
	if in == nil {
		in = new(EchoRequest)
	}
	out = new(EchoResponse)
	if err = c.Call("EchoService.EchoTwice", in, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DialEchoService connects to an EchoService at the specified network address.
func DialEchoService(network, addr string) (*EchoServiceClient, error) {
	c, err := protorpc.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &EchoServiceClient{c}, nil
}

// DialEchoServiceTimeout connects to an EchoService at the specified network address.
func DialEchoServiceTimeout(network, addr string, timeout time.Duration) (*EchoServiceClient, error) {
	c, err := protorpc.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	return &EchoServiceClient{c}, nil
}
