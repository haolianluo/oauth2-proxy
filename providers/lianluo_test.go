package providers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bmizerany/assert"
)

func testLianluoProvider(hostname string) *LianluoProvider {
	p := NewLianluoProvider(
		&ProviderData{
			ProviderName: "",
			LoginURL:     &url.URL{},
			RedeemURL:    &url.URL{},
			ProfileURL:   &url.URL{},
			ValidateURL:  &url.URL{},
			Scope:        ""})
	if hostname != "" {
		updateURL(p.Data().LoginURL, hostname)
		updateURL(p.Data().RedeemURL, hostname)
		updateURL(p.Data().ProfileURL, hostname)
		updateURL(p.Data().ValidateURL, hostname)
	}
	return p
}

func testLianluoBackend(payload string) *httptest.Server {
	path := "/account/v1/users"
	query := "access_token=imaginary_access_token&expand=verification"

	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			url := r.URL
			if url.Path != path || url.RawQuery != query {
				w.WriteHeader(404)
			} else {
				w.WriteHeader(200)
				w.Write([]byte(payload))
			}
		}))
}

func TestLianluoProviderDefaults(t *testing.T) {
	p := testLianluoProvider("")
	assert.NotEqual(t, nil, p)
	assert.Equal(t, "Lianluo", p.Data().ProviderName)
	assert.Equal(t, "https://mops-ucenter.lianluo.com/oauth2/auth",
		p.Data().LoginURL.String())
	assert.Equal(t, "https://mops-api.lianluo.com/account/v1/tokens",
		p.Data().RedeemURL.String())
	assert.Equal(t, "https://mops-api.lianluo.com/account/v1/users",
		p.Data().ValidateURL.String())
	assert.Equal(t, "", p.Data().Scope)
}

func TestLianluoProviderOverrides(t *testing.T) {
	p := NewLianluoProvider(
		&ProviderData{
			LoginURL: &url.URL{
				Scheme: "https",
				Host:   "mops-ucenter.lianluo.com",
				Path:   "/oauth2/auth"},
			RedeemURL: &url.URL{
				Scheme: "https",
				Host:   "mops-api.lianluo.com",
				Path:   "/account/v1/tokens"},
			ValidateURL: &url.URL{
				Scheme: "https",
				Host:   "mops-api.lianluo.com",
				Path:   "/account/v1/users"},
			Scope: ""})
	assert.NotEqual(t, nil, p)
	assert.Equal(t, "Lianluo", p.Data().ProviderName)
	assert.Equal(t, "https://mops-ucenter.lianluo.com/oauth2/auth",
		p.Data().LoginURL.String())
	assert.Equal(t, "https://mops-api.lianluo.com/account/v1/tokens",
		p.Data().RedeemURL.String())
	assert.Equal(t, "https://mops-api.lianluo.com/account/v1/users",
		p.Data().ValidateURL.String())
	assert.Equal(t, "", p.Data().Scope)
}

func TestLianluoProviderGetEmailAddress(t *testing.T) {
	b := testLianluoBackend("{\"email\": \"michael.bland@gsa.gov\"}")
	defer b.Close()

	b_url, _ := url.Parse(b.URL)
	p := testLianluoProvider(b_url.Host)

	session := &SessionState{AccessToken: "imaginary_access_token"}
	email, err := p.GetEmailAddress(session)
	assert.Equal(t, nil, err)
	assert.Equal(t, "michael.bland@gsa.gov", email)
}

// Note that trying to trigger the "failed building request" case is not
// practical, since the only way it can fail is if the URL fails to parse.
func TestLianluoProviderGetEmailAddressFailedRequest(t *testing.T) {
	b := testLianluoBackend("unused payload")
	defer b.Close()

	b_url, _ := url.Parse(b.URL)
	p := testLianluoProvider(b_url.Host)

	// We'll trigger a request failure by using an unexpected access
	// token. Alternatively, we could allow the parsing of the payload as
	// JSON to fail.
	session := &SessionState{AccessToken: "unexpected_access_token"}
	email, err := p.GetEmailAddress(session)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", email)
}

func TestLianluoProviderGetEmailAddressEmailNotPresentInPayload(t *testing.T) {
	b := testLianluoBackend("{\"foo\": \"bar\"}")
	defer b.Close()

	b_url, _ := url.Parse(b.URL)
	p := testLianluoProvider(b_url.Host)

	session := &SessionState{AccessToken: "imaginary_access_token"}
	email, err := p.GetEmailAddress(session)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", email)
}
