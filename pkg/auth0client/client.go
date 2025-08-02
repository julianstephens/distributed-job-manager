package auth0client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
)

type Auth0Client struct {
	conf         *models.Config
	client       *http.Client
	token        *string
	clientID     string
	clientSecret string
}

func NewAuth0Client(clientID string, clientSecret string) *Auth0Client {
	return &Auth0Client{
		conf:         config.GetConfig(),
		client:       &http.Client{},
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (a *Auth0Client) getToken() error {
	formData := url.Values{}
	formData.Add("grant_type", "client_credentials")
	formData.Add("client_id", a.clientID)
	formData.Add("client_secret", a.clientSecret)
	formData.Add("audience", a.conf.Auth0.Audience)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/oauth/token", a.conf.Auth0.Domain), strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var jsonData map[string]any
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return err
	}

	token := jsonData["access_token"].(string)
	if token != "" {
		a.token = &token
	}

	return nil
}

func (a *Auth0Client) Request(req *http.Request) (*http.Response, error) {
	if a.token == nil {
		a.getToken()
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *a.token))
	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		a.getToken()
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *a.token))
		return a.client.Do(req)
	}

	return res, nil
}
