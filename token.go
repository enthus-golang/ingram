package ingram

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   string    `json:"expires_in"`
	ValidUntil  time.Time `json:"-"`
}

func GetOAuthToken(ctx context.Context, clientID, clientSecret string) (*Token, error) {
	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	data.Add("client_id", clientID)
	data.Add("client_secret", clientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiEndpoint+"/oauth/oauth20/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	b, err := httputil.DumpRequest(req, true)
	fmt.Println(string(b))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	b, err = httputil.DumpResponse(res, true)
	fmt.Println(string(b))

	var t Token
	err = json.NewDecoder(res.Body).Decode(&t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (i *Ingram) checkAndUpdateToken(ctx context.Context) error {
	if i.token != nil && time.Now().Before(i.token.ValidUntil) {
		return nil
	}

	token, err := GetOAuthToken(ctx, i.clientID, i.clientSecret)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", token)

	expiresIn, err := strconv.Atoi(token.ExpiresIn)
	if err != nil {
		return err
	}
	token.ValidUntil = time.Now().Add(time.Duration(expiresIn-60) * time.Second)

	i.token = token

	return nil
}
