package token

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

func GetAccessToken(consumerKey, consumerSecret string) (string, error) {
	url := "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	authorizationString := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+authorizationString)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var accessTokenResponse AccessTokenResponse
	if err := json.NewDecoder(res.Body).Decode(&accessTokenResponse); err != nil {
		return "", err
	}

	return accessTokenResponse.AccessToken, nil
}
