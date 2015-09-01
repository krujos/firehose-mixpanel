package uaaclientcredentials

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//UAAClientCredentials provides a token for a given clientId and clientSecret.
//The token is refreshed for you according to expires_in
type UAAClientCredentials struct {
	uaaURI            *url.URL
	clientID          string
	clientSecret      string
	accessToken       string
	expiresAt         time.Time
	scope             string
	skipSSLValidation bool
}

//UAATokenResponse is the struct version of the json /oauth/token gives us
//when we ask for client credentials.
type UAATokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Jti         string `json:"jti"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

//GetBearerToken returns a currently valid bearer token to use against the
//CF API. You should not cache the token as the library will handle updating
//it if it's expired. This API will return an empty string and an error if
//there was a problem aquiring a token from UAA
func (creds *UAAClientCredentials) GetBearerToken() (string, error) {
	if time.Now().After(creds.expiresAt) {
		if err := creds.getToken(); nil != err {
			return "", err
		}
	}
	return "bearer " + creds.accessToken, nil
}

//New UAAClientCredentials factory
func New(uaaURI *url.URL, skipSSLValidation bool, clientID string,
	clientSecret string) (*UAAClientCredentials, error) {

	if len(clientID) < 1 {
		return nil, errors.New("clientID cannot be empty")
	}

	if len(clientSecret) < 1 {
		return nil, errors.New("clientSecret cannot be empty")
	}

	uri, _ := url.Parse(uaaURI.String() + "/oauth/token?grant_type=client_credentials")

	//Force the first call bo bearer token to get a new one.
	duration, _ := time.ParseDuration("-5m")
	expiresAt := time.Now().Add(duration)

	creds := &UAAClientCredentials{
		uaaURI:            uri,
		clientID:          clientID,
		clientSecret:      clientSecret,
		skipSSLValidation: skipSSLValidation,
		expiresAt:         expiresAt,
	}

	return creds, nil
}

func (creds *UAAClientCredentials) getTLSConfig() *tls.Config {
	if creds.skipSSLValidation {
		return &tls.Config{InsecureSkipVerify: true}
	}
	return &tls.Config{}
}

func (creds *UAAClientCredentials) getClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: creds.getTLSConfig(),
		},
	}
}

func (creds *UAAClientCredentials) getJSON() ([]byte, error) {
	client := creds.getClient()
	req, err := http.NewRequest("GET", creds.uaaURI.String(), nil)
	req.SetBasicAuth(creds.clientID, creds.clientSecret)

	resp, err := client.Do(req)
	if nil != err {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("UAA responded with bad status (" +
			strconv.Itoa(resp.StatusCode) + ")")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func (creds *UAAClientCredentials) getToken() error {

	body, err := creds.getJSON()

	var token UAATokenResponse
	json.Unmarshal(body, &token)

	if nil != err {
		return err
	}

	creds.accessToken = token.AccessToken
	//Give ourselves 1 min of buffer time for clock skews
	duration, _ := time.ParseDuration(strconv.Itoa(token.ExpiresIn-60) + "m")
	creds.expiresAt = time.Now().Add(duration)
	return nil

}
