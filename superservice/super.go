package superservice

import msg "github.com/playnb/mustang/msg.pb"

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	//"github.com/playnb/mustang/auth"
	"database/sql"
	"github.com/playnb/mustang/log"
	"github.com/playnb/mustang/utils"
)

//SuperService实例
var Instance = newSuperService()

func newSuperService() *superService {
	instance := new(superService)
	instance.rpc = newSuperRpc(instance)

	return instance
}

type superService struct {
	rpc       *superRpc
	db        *sql.DB
	findToken *sql.Stmt
	findSina  *sql.Stmt
}

func (service *superService) StartService(addr string) {
	go msg.ListenAndServeSuperService("tcp", addr, service.rpc)
	log.Trace("[SUPER] RPC在%s启动", addr)
	<-utils.CloseSig
}

func queryAccount(stms *sql.Stmt, key string) *msg.OAtuhUserProfile {
	rows, err := stms.Query(key)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	data := &msg.OAtuhUserProfile{}
	found := 0
	for rows.Next() {
		found++
		var accid uint64
		var nake_name string
		var token string
		err = rows.Scan(&accid, &nake_name, &token)
		data.Accid = proto.Uint64(accid)
		data.Name = proto.String(nake_name)
		data.Oauth = proto.String("SELF")
	}
	if found == 1 {
		return data
	} else {
		return nil
	}
}

func (service *superService) FindByToken(token string) *msg.OAtuhUserProfile {
	return queryAccount(service.findToken, token)
}

func (service *superService) FindBySinaID(sinaID string) *msg.OAtuhUserProfile {
	return queryAccount(service.findSina, sinaID)
}
