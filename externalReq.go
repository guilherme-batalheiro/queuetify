package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func basic_auth(client_id string, client_secret string) string {
	auth := client_id + ":" + client_secret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func request_current_song_info(token string) (map[string]interface{}, error) {
	var json_response map[string]interface{}
	var array_response []interface{}

	request, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 204 {
			return nil, errors.New("no song playing")
		}

		return nil, errors.New(resp.Status)
	}

	err = json.Unmarshal(responseBody, &json_response)
	if err != nil {
		return nil, err
	}

	progress_ms := json_response["progress_ms"].(float64)
	json_response = json_response["item"].(map[string]interface{})
	duration_ms := json_response["duration_ms"].(float64)
	song_name := json_response["name"].(string)

	array_response = json_response["artists"].([]interface{})
	artists := make([]string, len(array_response))
	for i, artist_json_info := range array_response {
		artists[i] = (artist_json_info.(map[string]interface{})["name"]).(string)
	}

	json_response = map[string]interface{}{
		"song_name":   song_name,
		"artists":     artists,
		"progress_ms": progress_ms,
		"duration_ms": duration_ms,
	}

	return json_response, nil
}

func request_user_info(token string) (map[string]interface{}, error) {
	var json_response map[string]interface{}

	request, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, err
	}

	// add query values
	query := request.URL.Query()
	query.Add("market", "PT")
	request.URL.RawQuery = query.Encode()

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	err = json.Unmarshal(responseBody, &json_response)
	if err != nil {
		return nil, err
	}

	id := json_response["id"].(string)
	display_name := json_response["display_name"].(string)
	email := json_response["email"].(string)

	json_response = map[string]interface{}{
		"id":           id,
		"display_name": display_name,
		"email":        email,
	}

	return json_response, nil
}

func request_refreshed_access_token(refresh_token string) (map[string]interface{}, error) {
	var json_response map[string]interface{}

	request, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", nil)
	if err != nil {
		return nil, err
	}

	// add query values
	query := request.URL.Query()
	query.Add("grant_type", "refresh_token")
	query.Add("refresh_token", refresh_token)
	request.URL.RawQuery = query.Encode()

	// add headers
	request.Header.Add("Authorization", "Basic "+basic_auth(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	err = json.Unmarshal(responseBody, &json_response)
	if err != nil {
		return nil, err
	}

	refresh_token = json_response["access_token"].(string)
	expires_in := json_response["expires_in"].(string)
	json_response = map[string]interface{}{
		"access_token": refresh_token,
		"expires_in":   expires_in,
	}

	return json_response, nil
}

func request_tokens(code string) (map[string]interface{}, error) {
	var json_response map[string]interface{}

	request, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", nil)
	if err != nil {
		return nil, err
	}

	// add query values
	query := request.URL.Query()
	query.Add("grant_type", "authorization_code")
	query.Add("code", code)
	query.Add("redirect_uri", "http://"+os.Getenv("ADDRESS")+"/")
	request.URL.RawQuery = query.Encode()

	// add headers
	request.Header.Add("Authorization", "Basic "+basic_auth(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	err = json.Unmarshal(responseBody, &json_response)
	if err != nil {
		return nil, err
	}

	access_token := json_response["access_token"].(string)
	refresh_token := json_response["refresh_token"].(string)
	expires_in := json_response["expires_in"].(float64)
	json_response = map[string]interface{}{
		"access_token":  access_token,
		"refresh_token": refresh_token,
		"expires_in":    expires_in,
	}

	return json_response, nil
}

func request_skip_song(token string) error {
	request, err := http.NewRequest(http.MethodPost, "https://api.spotify.com/v1/me/player/next", nil)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		defer resp.Body.Close()
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var error_message map[string]interface{}

		err = json.Unmarshal(responseBody, &error_message)
		if err != nil {
			return err
		}
		return errors.New(resp.Status)
	}

	return nil
}

func request_available_device_id(token string) (string, error) {
	var json_response map[string]interface{}

	request, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me/player/devices", nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	err = json.Unmarshal(responseBody, &json_response)
	if err != nil {
		return "", err
	}

	devices := json_response["devices"].([]interface{})
	if len(devices) == 0 {
		return "", errors.New("no devices available found")
	}

	first_device := devices[0].(map[string]interface{})
	device_id := first_device["id"].(string)

	return device_id, nil
}

func request_song_info(token string, song_name string) (map[string]interface{}, error) {
	var json_response map[string]interface{}
	var array_response []interface{}

	request, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/search", nil)
	if err != nil {
		return nil, err
	}

	// add query values
	query := request.URL.Query()
	query.Add("q", song_name)
	query.Add("type", "track")
	query.Add("market", "PT")
	query.Add("limit", "1")
	query.Add("offset", "0")
	request.URL.RawQuery = query.Encode()

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	err = json.Unmarshal(responseBody, &json_response)
	if err != nil {
		return nil, err
	}

	json_response = json_response["tracks"].(map[string]interface{})
	if json_response["total"].(float64) == 0 {
		return nil, errors.New("not found")
	}

	array_response = json_response["items"].([]interface{})
	json_response = array_response[0].(map[string]interface{})
	uri := json_response["uri"].(string)
	song_name = json_response["name"].(string)
	array_response = json_response["artists"].([]interface{})

	artists := make([]string, len(array_response))
	for i, artist_json_info := range array_response {
		artists[i] = (artist_json_info.(map[string]interface{})["name"]).(string)
	}

	json_response = map[string]interface{}{
		"uri":       uri,
		"song_name": song_name,
		"artists":   artists,
	}

	return json_response, nil
}

func request_add_music_to_queue_spotify(token string, device_id string, uri string) error {
	request, err := http.NewRequest(http.MethodPost, "https://api.spotify.com/v1/me/player/queue", nil)
	if err != nil {
		return err
	}

	// add query values
	query := request.URL.Query()
	query.Add("uri", uri)
	query.Add("device_id", device_id)
	request.URL.RawQuery = query.Encode()

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		defer resp.Body.Close()
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var error_message map[string]interface{}

		err = json.Unmarshal(responseBody, &error_message)
		if err != nil {
			return err
		}

		return errors.New(resp.Status)
	}

	return nil
}
