package spotify

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
