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

/*
Oauther is an interface, that if properly implemented, can be turned into a couple of Middleware that can capture an Oauth Token
*/
type Oauther interface {
	Config() *oauth.Config
	StartUrl() string
	RedirectUrl() string
}

// A Token handler allows you to specify what happens once you've received a token
type TokenHandler func(Oauther, *oauth.Token) rack.Middleware

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
	url := this.o.Config().AuthCodeURL(state)
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
		next()
		return
	}

	state1 := r.FormValue("state")
	state2, isstring := (sessioner.V)(vars).Clear("state").(string)

	//if states don't match, it's a potential CSRF attempt; we're just going to pass it on, and a 404 will probably be passed back (unless this happens to route somewhere else too)
	//perhaps we should just return a 401-Unauthorized, though
	if !isstring || state1 != state2 {
		//	Warning: Potential CSRF attempt : states don't match
		next()
		return
	}

	//Step 2: Exchange the code for the token
	code := r.FormValue("code")
	t := &oauth.Transport{Config: this.o.Config()}
	tok, _ := t.Exchange(code)

	//Step 3: Have some other middleware handle whatever they're doing with the token (probably logging a user in)
	process := this.t(this.o, tok)
	process.Run(vars, next)
}

//New converts your Oauther into Middleware, and allows you to do something once you have the token
func New(o Oauther, t TokenHandler) rack.Middleware {
	i := interceptor.New()
	i.Intercept(o.StartUrl(), &codeGetter{o})
	i.Intercept(o.RedirectUrl(), &tokenGetter{o, t})
	return i
}

//GetSite() is the function you use to actually get data from an Oauth provider
//It requires an Oauther, a token, a URL string, and a callback function for what to do with the response
func GetSite(o Oauther, tok *oauth.Token, site string, handler func(*http.Response) error) error {
	t := &oauth.Transport{Config: o.Config(), Token: tok}
	req, err := t.Client().Get(site)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	return handler(req)
}
