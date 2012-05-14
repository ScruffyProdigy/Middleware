/*
	the Facebooker package implements the Oauther interface, and provides facebook specific interactivity
*/
package facebooker

import (
	"github.com/HairyMezican/../goauth2/oauth"
	"strings"
)

type Data struct {
	// most of these variables should be straight up copied from https://developers.facebook.com/apps
	// if you mess up any of these variables, you will get an error
	// the fields have been tagged to make them easier to load in from a json file
	AppId       string   `json:"app_id"`       // on the dashboard - "App ID/API Key"
	AppSecret   string   `json:"app_secret"`   // on the dashboard - "App Secret"
	SiteUrl     string   `json:"site_url"`     // on the dashboard - "Site URL"
	Permissions []string `json:"permissions"`  // see http://developers.facebook.com/docs/authentication/permissions/ for more details
	StartUrl    string   `json:"start_url"`    // you decide - the route where the user should get start
	RedirectUrl string   `json:"redirect_url"` // you decide - where facebook sends the user after they've been authenticated
}

type Facebooker struct {
	data   Data
	config *oauth.Config
}

func New(data Data) *Facebooker {
	this := new(Facebooker)
	this.data = data
	return this
}

func (this *Facebooker) GetConfig() *oauth.Config {
	if this.config == nil {
		this.config = new(oauth.Config)
		this.config.ClientId = this.data.AppId
		this.config.ClientSecret = this.data.AppSecret
		this.config.Scope = strings.Join(this.data.Permissions, ",")
		this.config.AuthURL = "https://www.facebook.com/dialog/oauth"
		this.config.TokenURL = "https://graph.facebook.com/oauth/access_token"
		this.config.RedirectURL = this.data.SiteUrl + this.data.RedirectUrl
	}
	return this.config
}

func (this *Facebooker) GetStartUrl() string {
	return "/" + this.data.StartUrl
}

func (this *Facebooker) GetRedirectUrl() string {
	return "/" + this.data.RedirectUrl
}
