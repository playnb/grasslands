package superservice

import msg "github.com/playnb/mustang/msg.pb"

import (
	"github.com/golang/protobuf/proto"
	"github.com/playnb/mustang/auth"
	"github.com/playnb/mustang/log"
	"github.com/playnb/mustang/utils"
	//"strconv"
	"time"
)

type superRpc struct {
	super *superService //SuperService
	auths map[uint64]chan *auth.AuthUserData
}

func newSuperRpc(super *superService) *superRpc {
	ins := new(superRpc)
	ins.super = super
	ins.auths = make(map[uint64]chan *auth.AuthUserData)
	return ins
}

//心跳包
func (service *superRpc) Echo(req *msg.Ping, res *msg.Pong) error {
	res.Timestamp = proto.Uint32(uint32(time.Now().UnixNano() / int64(time.Millisecond)))
	res.ServiceId = proto.Uint64(req.GetServiceId())
	return nil
}

//登录SuperService
func (service *superRpc) Login(req *msg.LoginRequst, res *msg.LoginRespose) error {
	log.Trace("[SUPER] %s(%s) 的连接到来", utils.GetServiceName(int(req.GetServiceType())), string(req.GetServiceIp()))
	res.ServiceType = proto.Uint32(req.GetServiceType())
	res.ServiceId = proto.Uint64(100)
	res.ServiceIp = proto.String(req.GetServiceIp())
	res.RetCode = proto.Uint32(1)
	res.ExterPort = proto.Uint32(300)
	log.Trace("[SUPER] 返回 %s 服务器ID %d", utils.GetServiceName(int(res.GetServiceType())), res.GetServiceId())
	return nil
}

//获取SIna授权URL
func (service *superRpc) GetAuthUrlBySina(req *msg.OAuth2Request, res *msg.OAuth2Url) error {
	auth_sid := utils.GenEasyNextId(utils.SnowflakeSystemWork, utils.SnowflakeCatalogAuth)
	auth_ch := make(chan *auth.AuthUserData)
	service.auths[auth_sid] = auth_ch
	url, _ := auth.SinaOAuth2(auth_sid, func(data *auth.AuthUserData, accid uint64, err error) {
		if data != nil && err == nil {
			log.Trace("获得授权 auth_sid:%d data:%v", auth_sid, data)
			auth_ch <- data
		} else {
			log.Trace(err.Error())
			auth_ch <- nil
		}
	})
	res.Accid = proto.Uint64(req.GetAccid())
	res.AuthSid = proto.Uint64(auth_sid)
	res.Url = proto.String(url)
	return nil
}

//等待Sina授权返回结果
func (service *superRpc) WaitAuthResultBySina(req *msg.OAuth2Request, res *msg.OAuth2Response) error {
	if ch, ok := service.auths[req.GetAuthSid()]; ok {
		defer delete(service.auths, req.GetAuthSid())
		data := <-ch
		if data != nil {
			user := service.super.FindBySinaID(data.AuthUID)
			if user != nil {

			} else {
				res.RetCode = proto.Uint32(uint32(msg.OAtuhRetCode_AUTH_OK))
				res.Accid = proto.Uint64(req.GetAccid())
				res.AuthSid = proto.Uint64(req.GetAuthSid())
				res.User = new(msg.OAtuhUserProfile)
				res.User.Accid = proto.Uint64(req.GetAccid())
				res.User.Oauth = proto.String(data.ServiceName)
				res.User.Name = proto.String(data.AuthName)
			}
		} else {
			res.RetCode = proto.Uint32(uint32(msg.OAtuhRetCode_SINA_AUTH_FAILED))
		}
	}
	return nil
}

//Token登陆
func (service *superRpc) AuthByToken(req *msg.OAtuhTokenLogin, res *msg.OAuth2Response) error {
	data := service.super.FindBySinaID(req.GetToken())
	if data == nil {
		res.RetCode = proto.Uint32(uint32(msg.OAtuhRetCode_TOKEN_NOT_FOUND))
	} else {
		res.RetCode = proto.Uint32(uint32(msg.OAtuhRetCode_AUTH_OK))
		res.Accid = proto.Uint64(data.GetAccid())
		res.AuthSid = proto.Uint64(0)
		res.User = data
	}
	return nil
}
