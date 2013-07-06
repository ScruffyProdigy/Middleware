package googleplusser

import (
	"code.google.com/p/goauth2/oauth"
	"strings"
)

const (
	UserPermission  = "https://www.googleapis.com/auth/plus.me"
	EmailPermission = "https://www.googleapis.com/auth/userinfo.email"
)

type Data struct {
	// most of these variables should be straight up copied from https://code.google.com/apis/console/
	// if you mess up any of these variables, you will get an error
	ClientID     string `json:"client_id"`     // found in OAuth Client ID for Web Applications "Client ID"
	ClientSecret string `json:"client_secret"` // found in OAuth Client ID for Web Applications "Client secret"
	SiteUrl      string `json:"site_url"`      // found in OAuth Client ID for Web Applications "Redirect URIs" (the host of the specified URI)
	RedirectUri  string `json:"redirect_uri"`  // found in OAuth Client ID for Web Applications "Redirect URIs" (the path of the specified URI)
	Apikey       string `json:"api_key"`       // found in Simple API Access (Server Key) "API key"
	// these variables are yours to decide
	StartUri string `json:"start_uri"` // yours to decide, it is the path that you should direct the user to to log in
	//	these variables are created by us
	Permissions []string `json:"permissions"` // what you want to do, options found above "UserPermission is recommended"

}

type GooglePlus struct {
	data   Data
	config *oauth.Config
}

func New(data Data) *GooglePlus {
	gp := new(GooglePlus)
	gp.data = data
	return gp
}

func (this *GooglePlus) StartUrl() string {
	return "/" + this.data.StartUri
}

func (this *GooglePlus) RedirectUrl() string {
	return "/" + this.data.RedirectUri
}

func (this *GooglePlus) Config() *oauth.Config {
	if this.config == nil {
		this.config = new(oauth.Config)
		this.config.ClientId = this.data.ClientID
		this.config.ClientSecret = this.data.ClientSecret
		this.config.Scope = strings.Join(this.data.Permissions, ",")
		this.config.AuthURL = "https://accounts.google.com/o/oauth2/auth"
		this.config.TokenURL = "https://accounts.google.com/o/oauth2/token"
		this.config.RedirectURL = this.data.SiteUrl + this.data.RedirectUri
	}

	return this.config
}
