package test

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/playnb/mustang/global"
	"github.com/playnb/mustang/log"
)

func TestMysql() {
	db, err := sql.Open("mysql", global.MySqlUrl)
	if err != nil {
		log.Error(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Error(err.Error())
	}
	rows, err := db.Query("select ACCID,NAKE_NAME from ACCOUNT")
	if err != nil {
		log.Error(err.Error())
	}
	stmtIns, err := db.Prepare("insert into ACCOUNT (`ACCID`,`NAKE_NAME`) values( ?, ? )")
	if err != nil {
		log.Error(err.Error())
	}
	var result sql.Result
	result, err = stmtIns.Exec("666", "去")
	if err != nil {
		log.Error(err.Error())
	} else {
		n, _ := result.RowsAffected()
		log.Trace("影响的行数 %d", n)
	}
	stmtDel, err := db.Prepare("delete from `ACCOUNT` where `ACCID`=?")
	if err != nil {
		log.Error(err.Error())
	}
	stmtDel.Exec(666)
	defer rows.Close()
	for rows.Next() {
		var accid uint64
		var nake_name string
		err = rows.Scan(&accid, &nake_name)
		if err != nil {
			log.Error(err.Error())
		}
		log.Trace("找到数据 %d, %s", accid, nake_name)
	}
	////
	findToken, err := db.Prepare("select `ACCID`, `NAKE_NAME`,`TOKEN` from `ACCOUNT` where `TOKEN`=?")
	if err != nil {
		log.Error(err.Error())
	}
	rows, err = findToken.Query("EEE")

	findSina, err := db.Prepare("select  `ACCID`, `NAKE_NAME`,`TOKEN` from `ACCOUNT` where `SINA_UID`=?")
	rows, err = findSina.Query(1)

	log.Trace("=========================")
	if err != nil {
		log.Error(err.Error())
	}
	for rows.Next() {
		var accid uint64
		var nake_name string
		var token string
		err = rows.Scan(&accid, &nake_name, &token)
		if err != nil {
			log.Error(err.Error())
		}
		log.Trace("TOKEN找到数据 accis:%d, nake_name:%s token:%s", accid, nake_name, token)
	}
}
