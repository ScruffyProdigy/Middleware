# OAuther
This defines an interface for you to use for an OAuth provider, then takes that interface, and converts it into a Rack Middleware based system.  

## Dependencies
1.	It needs the session middleware to be inserted before it (github.com/HairyMezican/Middleware/session)
2.	It needs an interceptor to define it's routes (github.com/HairyMezican/Middleware/interceptor)
3.	It uses my adaptation of goauth2 to implement the oauth protocol (github.com/HairyMezican/goauth2/oauth)

## Installation
`go get github.com/HairyMezican/Middleware/oauther/...`

## Example

	package main

	import (
		"github.com/HairyMezican/Middleware/interceptor"
		"github.com/HairyMezican/Middleware/oauther"
		"github.com/HairyMezican/Middleware/oauther/facebooker"
		"github.com/HairyMezican/Middleware/sessioner"
		"github.com/HairyMezican/TheRack/rack"
		"github.com/HairyMezican/goauth2/oauth"
		"encoding/json"
		"net/http"
	)

	var data = facebooker.Data{
		AppId:       "123456789012345",                  //replace this with your own App ID from developers.facebook.com/apps
		AppSecret:   "1234567890abcdef0123456789abcdef", //replace this with your own App Secret from developers.facebook.com/apps
		SiteUrl:     "http://localhost:3000/",           // make sure this is the site url and port you specify at developers.facebook.com/apps
		Permissions: []string{},
		StartUrl:    "",
		RedirectUrl: "id",
	}

	func TokenHandler(o oauther.Oauther, tok *oauth.Token) rack.Middleware {
		fb := o.(*facebooker.Facebooker)
		return rack.Func(func(vars rack.Vars, next func()) {
			if tok == nil {
				rack.SetMessageString(vars, "User declined app")
			} else {
				rack.SetMessageString(vars, getUserID(fb, tok))
			}
		})
	}

	func getUserID(o oauther.Oauther, tok *oauth.Token) (result string) {
		oauther.GetSite(o, tok, "https://graph.facebook.com/me", func(res *http.Response) {
			//use json to read in the result, and get 
			var uid struct {
				ID string `json:"id"` //there are a lot of fields, but we really only care about the ID
			}

			d := json.NewDecoder(res.Body)
			err := d.Decode(&uid)
			if err != nil {
				panic(err)
			}

			result = uid.ID
		})
		return
	}

	func main() {
		cept := interceptor.New()

		fb := facebooker.New(data)
		oauther.SetIntercepts(cept, fb, TokenHandler)

		rackup := rack.New()
		rackup.Add(sessioner.Middleware)
		rackup.Add(cept)

		conn := rack.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
	

If you go to "localhost:3000/", you should be immediately redirected to facebook, and once you authorize the app, you'll be sent back, and you'll see your user ID