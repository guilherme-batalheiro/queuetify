package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func RequestUserInfo(accessToken string) (SpotifyUserInfo, error) {
	var json_response map[string]interface{}

	request, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return SpotifyUserInfo{}, err
	}

	// add query values
	query := request.URL.Query()
	query.Add("market", "PT")
	request.URL.RawQuery = query.Encode()

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return SpotifyUserInfo{}, err
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return SpotifyUserInfo{}, err
	}

	if resp.StatusCode != 200 {
		return SpotifyUserInfo{}, errors.New(string(responseBody))
	}

	err = json.Unmarshal(responseBody, &json_response)
	if err != nil {
		return SpotifyUserInfo{}, err
	}

	id := json_response["id"].(string)
	display_name := json_response["display_name"].(string)
	email := json_response["email"].(string)

	return SpotifyUserInfo{Id: id, DisplayName: display_name, Email: email}, nil
}

func requestAvailableDeviceId(accessToken string) (string, error) {
	var json_response map[string]interface{}

	request, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me/player/devices", nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
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

func requestSongInfo(accessToken string, song_name string) (SpotifySongInfo, error) {
	var json_response map[string]interface{}
	var array_response []interface{}

	request, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/search", nil)
	if err != nil {
		return SpotifySongInfo{}, err
	}

	// add query values
	query := request.URL.Query()
	query.Add("q", song_name)
	query.Add("type", "track")
	query.Add("market", "PT")
	query.Add("limit", "1")
	query.Add("offset", "0")
	request.URL.RawQuery = query.Encode()

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return SpotifySongInfo{}, err
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return SpotifySongInfo{}, err
	}

	if resp.StatusCode != 200 {
		return SpotifySongInfo{}, errors.New(resp.Status)
	}

	err = json.Unmarshal(responseBody, &json_response)
	if err != nil {
		return SpotifySongInfo{}, err
	}

	json_response = json_response["tracks"].(map[string]interface{})
	if json_response["total"].(float64) == 0 {
		return SpotifySongInfo{}, errors.New("not found")
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

	return SpotifySongInfo{Uri: uri, SongName: song_name, Artist: artists}, nil
}

func RequestAddMusicToQueue(accessToken string, song_query string) (SpotifySongInfo, error)  {
    if song_query == "" {
		return SpotifySongInfo{}, errors.New("The parameter song_query is empty.") 
    }


	request, err := http.NewRequest(http.MethodPost, "https://api.spotify.com/v1/me/player/queue", nil)
	if err != nil {
		return SpotifySongInfo{}, err
	}

	device_id, err := requestAvailableDeviceId(accessToken)
	if err != nil {
		return SpotifySongInfo{}, err
	}

    song_info, err := requestSongInfo(accessToken, song_query)
	if err != nil {
		return SpotifySongInfo{}, err
	}

	// add query values
	query := request.URL.Query()
	query.Add("uri", song_info.Uri)
	query.Add("device_id", device_id)
	request.URL.RawQuery = query.Encode()

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return SpotifySongInfo{}, err
	}

	if resp.StatusCode != 204 {
		defer resp.Body.Close()
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return SpotifySongInfo{}, err
		}

		var error_message map[string]interface{}

		err = json.Unmarshal(responseBody, &error_message)
		if err != nil {
			return SpotifySongInfo{}, err
		}

		return SpotifySongInfo{}, errors.New(resp.Status)
	}

	return song_info, nil 
}
