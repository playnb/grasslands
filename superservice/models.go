package superservice

import (
	_ "github.com/astaxie/beego/orm"
)

type Account struct {
	Accid          uint64
	NakeName       string
	Token          string
	SinaID         string
	last_logintime uint32
}

func init() {
	// 需要在init中注册定义的model
	//orm.RegisterModel(new(Account))
}
