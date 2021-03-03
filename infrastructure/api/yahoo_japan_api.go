package infrastructure

import (
	"atwell/authentication/usecase"
	"atwell/config"
	"net/http"
	"net/url"
)

type YahooJapanAuthAPI struct {
	Conf config.YahooAuthConfigurations
}

func NewYahooJapanAuthAPI(c config.YahooAuthConfigurations) usecase.YahooJapanAuthInfrastructure {
	return &YahooJapanAuthAPI{c}
}

func (a *YahooJapanAuthAPI) GetPublicKeyList() (*http.Response, error) {
	return http.Get("https://auth.login.yahoo.co.jp/yconnect/v2/public-keys")
}

func (a *YahooJapanAuthAPI) Token(code string) (*http.Response, error) {
	return http.PostForm("https://auth.login.yahoo.co.jp/yconnect/v2/token",
		url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {a.Conf.ClientID},
			"client_secret": {a.Conf.Secret},
			"redirect_uri":  {a.Conf.RedirectURL},
			"code":          {code},
		},
	)
}

func (a *YahooJapanAuthAPI) UserInfo(accessToken string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://userinfo.yahooapis.jp/yconnect/v1/attribute",
		nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	return client.Do(req)
}
