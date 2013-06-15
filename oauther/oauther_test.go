package oauther

import (
	"code.google.com/p/goauth2/oauth"
	"github.com/ScruffyProdigy/Middleware/interceptor"
	"github.com/ScruffyProdigy/Middleware/parser"
	"github.com/ScruffyProdigy/Middleware/redirecter"
	"github.com/ScruffyProdigy/Middleware/sessioner"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
)

func TokenHandlerFunc(o Oauther, tok *oauth.Token) rack.Middleware {
	return rack.Func(func(vars map[string]interface{}, next func()) {
		if tok == nil {
			(httper.V)(vars).SetMessageString("User declined app")
		} else {
			(httper.V)(vars).SetMessageString(getPayload(o, tok))
		}
	})
}

func getPayload(o Oauther, tok *oauth.Token) (result string) {
	GetSite(o, tok, "http://localhost:3001/data", func(res *http.Response) {
		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			result = err.Error()
		} else {
			result = string(bytes)
		}
	})
	return
}

type FakeProvider struct {
	MyURL       string
	YourURL     string
	YourLanding string
	config      *oauth.Config
	middleware  rack.Middleware
}

func (this FakeProvider) Config() *oauth.Config {
	if this.config == nil {
		this.config = new(oauth.Config)
		this.config.ClientId = "C713NT-1D"
		this.config.ClientSecret = "53CR3T"
		this.config.Scope = this.MyURL + "/whatever"
		this.config.AuthURL = this.MyURL + "/auth"
		this.config.TokenURL = this.MyURL + "/token"
		this.config.RedirectURL = this.YourURL + this.YourLanding
	}
	return this.config
}

func (FakeProvider) StartUrl() string {
	return "/start"
}

func (this FakeProvider) RedirectUrl() string {
	return this.YourLanding
}

func (this FakeProvider) Middleware() rack.Middleware {
	if this.middleware == nil {
		hostcept := interceptor.New()
		hostcept.Intercept("/auth", rack.Func(func(vars map[string]interface{}, next func()) {
			values := url.Values{}
			values.Set("state", parser.V(vars).FormValue("state"))
			values.Set("code", "c0D3")
			redirecter.V(vars).Redirect(parser.V(vars).FormValue("redirect_uri") + "?" + values.Encode())
		}))
		hostcept.Intercept("/token", rack.Func(func(vars map[string]interface{}, next func()) {
			if parser.V(vars).FormValue("code") == "c0D3" {
				httper.V(vars).GetRequest().Header.Set("content-type", "application/json")
				httper.V(vars).SetMessageString("{\"access_token\":\"tokendata\",\"refresh_token\":\"refreshtoken1\",\"expires_in\":3600}")
			}
		}))
		hostcept.Intercept("/data", rack.Func(func(vars map[string]interface{}, next func()) {
			if auth := httper.V(vars).GetRequest().Header.Get("Authorization"); auth != "Bearer tokendata" {
				httper.V(vars).SetMessageString("Invalid authorization: " + auth)
			}
			httper.V(vars).SetMessageString("payload")
		}))
		hostrackup := rack.New()
		hostrackup.Add(sessioner.Middleware)
		hostrackup.Add(hostcept)

		this.middleware = hostrackup
	}
	return this.middleware
}

func init() {
	//set up oauth host

	hoster := &FakeProvider{
		YourURL:     "http://localhost:3000",
		YourLanding: "/callback",
		MyURL:       "http://localhost:3001",
	}

	hostconn := httper.HttpConnection(":3001")
	go hostconn.Go(hoster.Middleware())

	//set up our site
	cept := interceptor.New()
	SetIntercepts(cept, hoster, TokenHandlerFunc)

	rackup := rack.New()
	rackup.Add(sessioner.Middleware)
	rackup.Add(cept)

	conn := httper.HttpConnection(":3000")
	go conn.Go(rackup)
}

func GetFrom(url string, cookies []*http.Cookie, t *testing.T) (int, string, []*http.Cookie) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if cookies != nil {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	client := &http.Client{Jar: jar}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer res.Body.Close()

	text, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	return res.StatusCode, string(text), res.Cookies()
}

func GetFromResults(url string, cookies *[]*http.Cookie, t *testing.T, assertedStatus int, assertedText string) {
	var status int
	var text string
	status, text, *cookies = GetFrom(url, *cookies, t)
	if status != assertedStatus {
		t.Errorf("Status should be %d, not %d", assertedStatus, status)
	}
	if text != assertedText {
		t.Errorf("Text should be %s, not %s", assertedText, text)
	}
}

func Test_CSRF(t *testing.T) {
	var cookies []*http.Cookie
	GetFromResults("http://localhost:3000/callback", &cookies, t, 404, "")
}

func Test_Normal(t *testing.T) {
	var cookies []*http.Cookie
	GetFromResults("http://localhost:3000/start", &cookies, t, 200, "payload")
}
