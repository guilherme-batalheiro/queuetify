package models

import "sync"

type User struct {
	SpotifyId   string
	DisplayName string
	Email       string

	AccessToken    string
	AccessTokenExp float64
	RefreshToken   string

	RoomCode string
}

type UserModel struct {
	DB *sync.Map
}

func (m *UserModel) Insert(spotify_id string, display_name string, email string) (bool, error) {

	m.DB.Store(spotify_id, &User{SpotifyId: spotify_id, DisplayName: display_name, Email: email, RoomCode: ""})

	return true, nil
}

func (m *UserModel) Get(spotify_id string) (User, bool) {

	user, ok := m.DB.Load(spotify_id)

	if !ok {
		return User{}, false
	}

	return *user.(*User), ok
}

func (m *UserModel) UpdateTokens(spotify_id string, accessToken string, accessTokenExp float64, refreshToken string) bool {

	user, ok := m.DB.Load(spotify_id)

	if !ok {
		return false
	}

	user.(*User).AccessToken = accessToken
	user.(*User).AccessTokenExp = accessTokenExp
	user.(*User).RefreshToken = refreshToken

	return true
}

func (m *UserModel) GetRoomCode(spotify_id string) (string, bool) {

	user, ok := m.DB.Load(spotify_id)

	if !ok || user.(*User).RoomCode == "" {
		return "", false
	}

	return user.(*User).RoomCode, true
}

func (m *UserModel) UpdateRoomCode(spotify_id string, roomCode string) bool {

	user, ok := m.DB.Load(spotify_id)

	if !ok {
		return false
	}

	user.(*User).RoomCode = roomCode

	return true
}
