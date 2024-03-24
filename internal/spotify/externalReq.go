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
