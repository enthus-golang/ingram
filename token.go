package ingram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
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

func (i *Ingram) GetOAuthToken(ctx context.Context, clientID, clientSecret string) (*Token, error) {
	data := fmt.Sprintf(`grant_type=client_credentials&client_id=%s&client_secret=%s`, clientID, clientSecret)

	version := "30"
	if i.isSandbox {
		version = "20"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/oauth/oauth%s/token", apiEndpoint, version), strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	if i.logger != nil {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		i.logger.Printf(string(b))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if i.logger != nil {
		b, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, err
		}
		i.logger.Printf(string(b))
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("unable to create token")
	}

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

	token, err := i.GetOAuthToken(ctx, i.clientID, i.clientSecret)
	if err != nil {
		return err
	}

	expiresIn, err := strconv.Atoi(token.ExpiresIn)
	if err != nil {
		return err
	}
	token.ValidUntil = time.Now().Add(time.Duration(expiresIn-60) * time.Second)

	i.token = token

	return nil
}
