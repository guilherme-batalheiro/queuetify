package models

import "sync"

type Room struct {
	Code           string
	OwnerSpotifyId string
}

type RoomModel struct {
	DB *sync.Map
}

func (m *RoomModel) Insert(roomCode string, ownerSpotifyId string) (bool, error) {

	m.DB.Store(roomCode, &Room{Code: roomCode, OwnerSpotifyId: ownerSpotifyId})

	return true, nil
}

func (m *RoomModel) Get(roomCode string) (Room, bool) {
	room, ok := m.DB.Load(roomCode)

	if !ok {
		return Room{}, false
	}

	return *room.(*Room), ok
}

func (m *RoomModel) Delete(roomCode string) {
	m.DB.Delete(roomCode)
}
