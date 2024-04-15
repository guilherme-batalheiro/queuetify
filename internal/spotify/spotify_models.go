package spotify

type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    float64
}

type SpotifyUserInfo struct {
	Id          string
	DisplayName string
	Email       string
}

type SpotifySongInfo struct {
	Uri      string
	SongName string
	Artist   []string
}
