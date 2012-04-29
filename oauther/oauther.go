/*
	Oather provides an interface for any OAuth service to work in with a rack based system
*/

package oauther

import (
	"github.com/HairyMezican/goauth2/oauth"
	"github.com/HairyMezican/TheRack/rack"
	"github.com/HairyMezican/Middleware/interceptor"
	"github.com/HairyMezican/Middleware/sessioner"
	"net/http"
	"crypto/rand"
	"encoding/base64"
)

type TokenHandler func(Oauther,*oauth.Token) rack.Middleware

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

func (this codeGetter) Run(r *http.Request, vars rack.Vars, next rack.Next) (status int, header http.Header, message []byte) {
	state := randomString()
	vars.Apply(sessioner.Set("state", state))
	w := rack.BlankResponse()
	url := this.o.GetConfig().AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
	return w.Results()
}

type tokenGetter struct {
	o Oauther
	t TokenHandler
}

func (this tokenGetter) Run(r *http.Request, vars rack.Vars, next rack.Next) (status int, header http.Header, message []byte) {
	//Step 1: Ensure states match
	state1 := r.FormValue("state")
	state2, isString := vars.Apply(sessioner.Clear("state")).(string)

	//if states don't match, it's a potential CSRF attempt; we're just going to pass it on, and a 404 will probably be passed back (unless this happens to route somewhere else too)
	//perhaps we should just return a 401-Unauthorized, though
	if !isString {
		//	Warning: Potential CSRF attempt : cookie not set properly
		return next()
	}
	if state1 != state2 {
		//	Warning: Potential CSRF attempt : states don't match
		return next()
	}

	//Step 2: Exchange the code for the token
	code := r.FormValue("code")
	t := &oauth.Transport{oauth.Config: this.o.GetConfig()}
	tok, _ := t.Exchange(code)

	//Step 3: Have some other middleware handle whatever they're doing with the token (probably logging a user in)
	process := this.t(this.o,tok)
	return process.Run(r, vars, next)
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
