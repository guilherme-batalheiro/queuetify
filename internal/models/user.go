package models

import (
	"sync"

	"queuetify.gbatalheiro.pt/internal/spotify"
)

type User struct {
	SpotifyId   string
	DisplayName string
	Email       string

	Tokens spotify.Tokens

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

func (m *UserModel) UpdateTokens(spotify_id string, tokens spotify.Tokens) bool {

	user, ok := m.DB.Load(spotify_id)

	if !ok {
		return false
	}

    user.(*User).Tokens = tokens

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
