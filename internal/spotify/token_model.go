package spotify

import (
	"time"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	AuthHeader   string
}

func (t *Tokens) GetAccessToken() (string, error) {
	currentTime := time.Now()
	unixSeconds := currentTime.Unix()

	if unixSeconds < t.ExpiresIn {
		return t.AccessToken, nil
	}

	newTokens, err := GetSpotifyAuthTokensFromRefresh(*t)
    if err != nil {
        return "", err
    }

	t.AccessToken = newTokens.AccessToken
	t.ExpiresIn = newTokens.ExpiresIn

	return t.AccessToken, nil
}
