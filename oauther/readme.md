# OAuther
This defines an interface for you to use for an OAuth provider, then takes that interface, and converts it into a Rack Middleware based system.  

## Dependencies
It needs the sessioner middleware to be inserted before it (github.com/ScruffyProdigy/Middleware/sessioner)

## Documentation
http://godoc.org/github.com/ScruffyProdigy/Middleware/oauther

## Installation
`go get github.com/ScruffyProdigy/Middleware/oauther/...`

## Usage

* Find or create an implementation of Oauther for the site you wish to access
	* Google+ and Facebook implementations can be found in this project
* Fill out the needed fields within your implementation
	* Typically, a description of each of the fields and where to find them are provided
* Call oauther.New() to convert your Oauther into a Middleware that you should add to your rack
	* The first parameter is the Oauther you filled out in the previous step
	* The second parameter is the function you will create in the next step
* A token is needed to access info on the website you're contacting; create the function that describes what to do once you have obtained the token
	* The function should take your Oauther, and the token obtained, and return a middleware
		* Generally, the middleware will need to store the token somewhere (your call on this)
		* And/or the middleware will need to immediately access the website (see the next step)
	* You should add the resulting Middleware to your rack somewhere after sessioner
* Any time you want to access information from the site, call GetSite()
	* The first two parameters are the Oauther and the token
	* The third parameter is the URL you want to access information from
	* The fourth parameter is a function that handles the response received from the website

## Example

	package main

	import (
		"github.com/ScruffyProdigy/Middleware/oauther"
		"github.com/ScruffyProdigy/Middleware/oauther/facebooker"
		"github.com/ScruffyProdigy/Middleware/sessioner"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
		"code.google.com/p/goauth2/oauth"
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
		return rack.Func(func(vars map[string]interface{}, next func()) {
			if tok == nil {
				(httper.V)(vars).SetMessageString("User declined app")
			} else {
				(httper.V)(vars).SetMessageString(getUserID(o, tok))
			}
		})
	}

	func getUserID(o oauther.Oauther, tok *oauth.Token) (result string) {
		err := oauther.GetSite(o, tok, "https://graph.facebook.com/me", func(res *http.Response) error {
			//use json to read in the result, and get 
			var uid struct {
				ID string `json:"id"` //there are a lot of fields, but we really only care about the ID
			}

			d := json.NewDecoder(res.Body)
			err := d.Decode(&uid)
			if err != nil {
				return err
			}

			result = uid.ID
			return nil
		})
		if err != nil {
			result = err.Error()
		}
		return
	}

	func main() {
		fb := facebooker.New(data)
		cept := oauther.New(fb, TokenHandler)

		rackup := rack.New()
		rackup.Add(sessioner.Middleware)
		rackup.Add(cept)

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}

If you go to "localhost:3000/", you should be immediately redirected to facebook, and once you authorize the app, you'll be sent back, and you'll see your user ID