package main

import (
	//"TestGo/Msg"
	"github.com/playnb/mustang/auth"
	"github.com/playnb/mustang/weixin"
	//"github.com/playnb/mustang/gate"
	"github.com/playnb/mustang/log"
	//"github.com/playnb/mustang/network"
	//"TestGo/mustang/network/protobuf"
	"github.com/playnb/mustang/utils"
	//"github.com/golang/protobuf/proto"
	"github.com/playnb/grasslands/gateservice"
	"github.com/playnb/grasslands/superservice"
	"github.com/playnb/grasslands/test"
	"github.com/playnb/mustang/global"
	"github.com/playnb/mustang/nosql"
	"github.com/playnb/mustang/sqldb"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

/*
oqxYYs5Y8h0_UTpKJLngn6VLNb_8
oqxYYs-koxpfISmvqFJK-Z3DveAc
*/

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}
func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

var baseDir = getParentDirectory(getCurrentDirectory()) + "/src/github.com/playnb/grasslands/"

func handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	index := strings.LastIndex(path, ".")
	if index != -1 {
		request_type := path[index:]
		switch request_type {
		case ".css":
			w.Header().Set("content-type", "text/css")
		case ".js":
			w.Header().Set("content-type", "text/javascript")
		default:
		}
	}
	path = baseDir + path
	log.Trace("获取静态文件: %s", path)
	fin, err := os.Open(path)
	defer fin.Close()
	if err != nil {
		log.Error("static resource: %v", err)
		w.WriteHeader(505)
	} else {
		fd, _ := ioutil.ReadAll(fin)
		w.Write(fd)
	}
}

func handleAppDemo(w http.ResponseWriter, r *http.Request) {
	log.Debug("Visit: %s", r.Proto+r.Host+r.RequestURI)
	r.ParseForm()
	var data struct {
		AccessToken string
		OpenID      string
		TimeStamp   string
		Nonce       string
		Signature   string
		ShareOpenID string
	}
	data.AccessToken = strings.Join(r.Form["access_token"], "")
	data.OpenID = strings.Join(r.Form["openid"], "")
	data.TimeStamp = strings.Join(r.Form["timestamp"], "")
	data.Nonce = strings.Join(r.Form["nonce"], "")
	data.ShareOpenID = strings.Join(r.Form["shareOpenID"], "")
	data.Signature = weixin.MakeSignature_js(data.TimeStamp, data.Nonce, weixin.Profile().JsApiTicket, "http://"+r.Host+r.RequestURI)
	path := baseDir + "/htdoc/app/demo.html"
	tpl, err := template.ParseFiles(path)
	if err != nil {
		log.Debug(err.Error())
		return
	}
	err = tpl.Execute(w, data)
	if err != nil {
		log.Debug(err.Error())
		return
	}
}

//作为一个代理访问微信web的api(等弄明白js咋个跨域就不用这个了)
func handleWxApi(w http.ResponseWriter, r *http.Request) {
	log.Debug("Visit: %s", r.RequestURI)
	url := r.RequestURI
	url = strings.Replace(url, "/wx/api/", "https://api.weixin.qq.com/", 1)
	log.Debug("API: %s", url)
	res, err := http.Get(url)
	if err != nil {
		log.Debug(err.Error())
		return
	}
	result, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	log.Debug("API_RET: %s", string(result))
	w.Write(result)
}

func UserAuthorized(state string, token *weixin.WebAccessToken, w http.ResponseWriter, r *http.Request) bool {
	key := "NikeName:" + token.OpenID
	ret, err := nosql.Redis.Get(key).Result()
	if err != nil {
		userinfo := weixin.GetUserInfo(token.OpenID, true)
		if userinfo != nil {
			ret, err = nosql.Redis.Set(key, userinfo.NikeName, 0).Result()
			if err != nil {
				log.Debug(err.Error() + " | " + ret)
			} else {
				log.Trace("App: %s 有人授权了(new)name: %s", state, userinfo.NikeName)
			}
		}
	} else {
		log.Trace("App: %s 有人授权了name: %s", state, ret)
	}
	return true
}

func main() {
	nosql.InitRedis(global.RedisURL, global.RedisPin, 0)
	sqldb.Init()
	mux := http.NewServeMux()
	weixin.InitWeiXin(&weixin.WeiXinConfig{AppID: global.AppID,
		AppSecret:       global.AppSecret,
		OwnerID:         global.OwnerID,
		ProcessMsg:      nil,
		UserAuthorized:  UserAuthorized,
		ServiceToken:    "test_wx",
		MyServiceDomain: global.ServiceDomain,
		AuthRedirectURL: "http://" + global.ServiceDomain + "/wx/authorize_redirect"}, mux)
	weixin.RegistAuthRedirectUrl("FIRST_TEST", "http://"+global.ServiceDomain+"/app/demo")
	mux.HandleFunc("/htdoc/", handleStatic)
	mux.HandleFunc("/app/demo", handleAppDemo)
	mux.HandleFunc("/wx/api/", handleWxApi)

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Error(err.Error())
	}
	//if err := http.ListenAndServeTLS(":3000", baseDir+"https/cert.pem", baseDir+"https/key.pem", mux); err != nil {
	//	log.Error(err.Error())
	//}
	return

	log.Trace("启动......")

	go superservice.Instance.StartService(utils.SuperRpcAddr)
	go gateservice.Instance.StartService(utils.SuperRpcAddr)

	<-utils.CloseSig
	log.Trace("结束......")
	return

	test.TestSinaAuth()
	auth.InitAuthHttpService()
	test.TestProtorpc()
	utils.TestSnow()
}
