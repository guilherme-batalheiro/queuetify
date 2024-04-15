package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

func GenerateAuthLink(clientId string, redirectUrl string) (string, error) {
	baseURL := "http://accounts.spotify.com/authorize"

	values := url.Values{}
	values.Add("response_type", "code")
	values.Add("client_id", clientId)
	values.Add("scope", "user-read-email user-modify-playback-state user-read-playback-state user-read-currently-playing")
	values.Add("redirect_uri", "http://"+redirectUrl+"/userSpotifyLoginCallback")

	return baseURL + "?" + values.Encode(), nil
}

func basicAuth(client_id string, client_secret string) string {
	auth := client_id + ":" + client_secret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func GetSpotifyAuthTokens(authCode string, clientId string, clientSecret string, redirectUrl string) (Tokens, error) {
	var json_response map[string]interface{}

	request, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", nil)
	if err != nil {
		return Tokens{}, err
	}

	// add query values
	query := request.URL.Query()
	query.Add("grant_type", "authorization_code")
	query.Add("code", authCode)
	query.Add("redirect_uri", "http://"+redirectUrl+"/userSpotifyLoginCallback")
	request.URL.RawQuery = query.Encode()

	// add headers
	auth := "Basic " + basicAuth(clientId, clientSecret)
	request.Header.Add("Authorization", auth)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return Tokens{}, err
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Tokens{}, err
	}

	println(clientId)
	if resp.StatusCode != 200 {
		return Tokens{}, errors.New(string(responseBody))
	}

	err = json.Unmarshal(responseBody, &json_response)
	if err != nil {
		return Tokens{}, err
	}

	accessToken := json_response["access_token"].(string)
	refreshToken := json_response["refresh_token"].(string)

	currentTime := time.Now()
	unixSeconds := currentTime.Unix()
	expiresIn := unixSeconds + int64(json_response["expires_in"].(float64))

	return Tokens{AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: expiresIn, AuthHeader: auth}, nil
}

func GetSpotifyAuthTokensFromRefresh(token Tokens) (Tokens, error) {
	var json_response map[string]interface{}

	request, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", nil)
	if err != nil {
		return Tokens{}, err
	}

	// add query values
	query := request.URL.Query()
	query.Add("grant_type", "refresh_token")
	query.Add("refresh_token", token.RefreshToken)
	request.URL.RawQuery = query.Encode()

	// add headers
	request.Header.Add("Authorization", token.AuthHeader)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return Tokens{}, err
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Tokens{}, err
	}

	if resp.StatusCode != 200 {
		return Tokens{}, errors.New(resp.Status)
	}

	err = json.Unmarshal(responseBody, &json_response)
	if err != nil {
		return Tokens{}, err
	}

	accessToken := json_response["access_token"].(string)
	currentTime := time.Now()
	unixSeconds := currentTime.Unix()
	expiresIn := unixSeconds + int64(json_response["expires_in"].(float64))

	return Tokens{AccessToken: accessToken, RefreshToken: token.RefreshToken, ExpiresIn: expiresIn, AuthHeader: token.AuthHeader}, nil
}
