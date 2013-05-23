/*
	Oather provides an interface for any OAuth service to work in with a rack based system
*/

package oauther

import (
	"code.google.com/p/goauth2/oauth"
	"crypto/rand"
	"encoding/base64"
	"github.com/ScruffyProdigy/Middleware/interceptor"
	"github.com/ScruffyProdigy/Middleware/redirecter"
	"github.com/ScruffyProdigy/Middleware/sessioner"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"net/http"
)

type TokenHandler func(Oauther, *oauth.Token) rack.Middleware

type Oauther interface {
	GetConfig() *oauth.Config
	GetStartUrl() string
	GetRedirectUrl() string
}

type codeGetter struct {
	o Oauther
}

func randomString() string {
	b := make([]byte, 80)
	rand.Read(b)
	en := base64.StdEncoding
	d := make([]byte, en.EncodedLen(len(b)))
	en.Encode(d, b)
	return string(d)
}

func (this codeGetter) Run(vars map[string]interface{}, next func()) {
	state := randomString()
	(sessioner.V)(vars).Set("state", state)
	url := this.o.GetConfig().AuthCodeURL(state)
	(redirecter.V)(vars).Redirect(url)
}

type tokenGetter struct {
	o Oauther
	t TokenHandler
}

func (this tokenGetter) Run(vars map[string]interface{}, next func()) {
	//Step 1: Ensure states match
	r := httper.V(vars).GetRequest()
	if r == nil {
		panic("Request not found")
	}

	state1 := r.FormValue("state")
	state2 := (sessioner.V)(vars).Clear("state")

	//if states don't match, it's a potential CSRF attempt; we're just going to pass it on, and a 404 will probably be passed back (unless this happens to route somewhere else too)
	//perhaps we should just return a 401-Unauthorized, though
	if state1 != state2 {
		//	Warning: Potential CSRF attempt : states don't match
		next()
		return
	}

	//Step 2: Exchange the code for the token
	code := r.FormValue("code")
	t := &oauth.Transport{oauth.Config: this.o.GetConfig()}
	tok, _ := t.Exchange(code)

	//Step 3: Have some other middleware handle whatever they're doing with the token (probably logging a user in)
	process := this.t(this.o, tok)
	process.Run(vars, next)
}

func SetIntercepts(i interceptor.Interceptor, o Oauther, t TokenHandler) {
	i.Intercept(o.GetStartUrl(), &codeGetter{o})
	i.Intercept(o.GetRedirectUrl(), &tokenGetter{o, t})
}

func GetSite(o Oauther, tok *oauth.Token, site string, handler func(*http.Response)) {
	t := &oauth.Transport{oauth.Config: o.GetConfig(), oauth.Token: tok}
	req, err := t.Client().Get(site)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	handler(req)
}
