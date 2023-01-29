package authentication

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rte-etl-routine/config"
	"sync"
	"time"
)

type AuthResponse struct {
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   time.Duration `json:"expires_in"`
	ExpiresAt   time.Time
}

type Authenticator struct{}

const TokenExpirationSafetyDelta time.Duration = 10

var (
	authResponse *AuthResponse
	mutex        sync.Mutex
)

func New() *Authenticator {
	return &Authenticator{}
}

func (a *Authenticator) authenticate() error {
	client := &http.Client{}
	request, err := http.NewRequest("POST", config.GetEnv().Auth.URI, nil)
	request.Header.Set("Authorization", "Basic "+config.GetEnv().Auth.Code)
	token, err := client.Do(request)
	if err != nil {
		fmt.Println("Unable to authenticate", err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Unable to close Body reader ", err)
		}
	}(token.Body)

	body, err := io.ReadAll(token.Body)
	if err != nil {
		fmt.Println("Unable to read auth response body")
	}

	var response AuthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Can not unmarshal JSON from authentication")
		return err
	}

	response.ExpiresAt = time.Now().Add(time.Second * response.ExpiresIn)
	authResponse = &response
	return nil
}

func (a *Authenticator) GetToken() (string, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if authResponse == nil || isTokenExpired() {
		err := a.authenticate()
		if err != nil {
			return "", err
		}
	}
	return authResponse.AccessToken, nil
}

func isTokenExpired() bool {
	return authResponse.ExpiresAt.Before(time.Now().Add(-time.Second * TokenExpirationSafetyDelta))
}
