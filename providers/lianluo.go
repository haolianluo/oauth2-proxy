package providers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/bitly/oauth2_proxy/api"
)

type LianluoProvider struct {
	*ProviderData
}

func NewLianluoProvider(p *ProviderData) *LianluoProvider {
	p.ProviderName = "GitLab"
	if p.LoginURL == nil || p.LoginURL.String() == "" {
		p.LoginURL = &url.URL{
			Scheme: "https",
			Host:   "mops-ucenter.lianluo.com",
			Path:   "/oauth2/auth",
		}
	}
	if p.RedeemURL == nil || p.RedeemURL.String() == "" {
		p.RedeemURL = &url.URL{
			Scheme: "https",
			Host:   "mops-api.lianluo.com",
			Path:   "/account/v1/tokens",
		}
	}
	if p.ValidateURL == nil || p.ValidateURL.String() == "" {
		p.ValidateURL = &url.URL{
			Scheme: "https",
			Host:   "mops-api.lianluo.com",
			Path:   "/account/v1/users",
		}
	}
	if p.Scope == "" {
		p.Scope = ""
	}
	return &LianluoProvider{ProviderData: p}
}

func (p *LianluoProvider) GetEmailAddress(s *SessionState) (string, error) {

	req, err := http.NewRequest("GET",
		p.ValidateURL.String()+"?access_token="+s.AccessToken, nil)
	if err != nil {
		log.Printf("failed building request %s", err)
		return "", err
	}
	json, err := api.Request(req)
	if err != nil {
		log.Printf("failed making request %s", err)
		return "", err
	}
	return json.Get("email").String()
}
