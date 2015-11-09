package test

import (
	"github.com/playnb/mustang/auth"
	"github.com/playnb/mustang/log"
	"net/http"
	"text/template"
)

var notAuthenticatedTemplate = template.Must(template.New("").Parse(`
<html><body>
<form action="/sina/authorize" method="POST"><input type="submit" value="测试新浪微博授权"/></form>
</form>
</body></html>
`))

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Trace("\n--")
	log.Trace("Visit: " + r.RequestURI)
	notAuthenticatedTemplate.Execute(w, nil)
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {
	log.Trace("\n--")
	log.Trace("Visit: " + r.RequestURI)
	//Get the Google URL which shows the Authentication page to the user
	url, _ := auth.SinaOAuth2(123456, func(data *auth.AuthUserData, accid uint64, err error) {
		if data != nil && err == nil {
			log.Trace("获得授权 %d", accid)
			log.Trace("%v", data)
			http.Redirect(data.Response, data.Request, "http://www.tuntun.site/", http.StatusFound)
		} else {
			log.Trace(err.Error())
		}
	})
	log.Trace("URL: %v\n", url)
	//redirect user to that page
	http.Redirect(w, r, url, http.StatusFound)
}

func TestSinaAuth() {
	mux := http.NewServeMux()
	mux.HandleFunc("/sina", handleRoot)
	mux.HandleFunc("/sina/authorize", handleAuthorize)
	go http.ListenAndServe(":3000", mux)
	log.Debug("测试连接在 127.0.0.1:3000")
}
