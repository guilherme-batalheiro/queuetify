package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var users_number int = 0
var users_map map[string]int

func generate_room_code() string {

	n, err := strconv.Atoi(os.Getenv("CODE_SIZE"))
	if err != nil {
		log.Fatal(err)
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func get_update_access_token(user_id string) (string, error) {
	refresh_token, err := get_refresh_token_from_db(user_id)
	if err != nil {
		return "", fmt.Errorf("Failed to get refresh token from db error: %w", err)
	}

	response_tokens, err := request_refreshed_access_token(refresh_token)
	if err != nil {
		return "", fmt.Errorf("Failed to request refreshed acess token: %w", err)
	}

	err = update_access_token_in_db(user_id, response_tokens)
	if err != nil {
		return "", fmt.Errorf("Failed to update access token in data base: %w", err)
	}

	return response_tokens["access_token"].(string), nil
}

func get_access_token(user_id string) (string, error) {
	var acc_token string

	acc_token_is_valid, err := access_token_is_valid_db(user_id)
	if err != nil {
		return "", fmt.Errorf("Failed to check if token is valid in database: %w", err)
	}

	if acc_token_is_valid {
		acc_token, err = get_access_token_from_db(user_id)
		if err != nil {
			return "", fmt.Errorf("Failed to get access token from db: %w", err)
		}
	} else {
		acc_token, err = get_update_access_token(user_id)
		if err != nil {
			return "", fmt.Errorf("Failed to get update access token: %w", err)
		}
	}

	return acc_token, nil
}

func auth_func(code string) (map[string]interface{}, error) {
	tokens, err := request_tokens(code)
	if err != nil {
		return nil, fmt.Errorf("Failed to request tokens: %w", err)
	}

	acc_token := tokens["access_token"].(string)
	user_info, err := request_user_info(acc_token)
	if err != nil {
		return nil, fmt.Errorf("Failed to request user info: %w", err)
	}

	user_id := user_info["id"].(string)
	exits, err := check_if_user_exits_in_db(user_id)
	if err != nil {
		return nil, fmt.Errorf("Failed to check if user exits in data base: %w", err)
	}

	if !exits {
		err = add_user_to_db(tokens, user_info)
		if err != nil {
			return nil, fmt.Errorf("Failed to add user to data base: %w", err)
		}
	} else {
		err = update_user_in_db(tokens, user_info)
		if err != nil {
			return nil, fmt.Errorf("Failed to update user in data base: %w", err)
		}
	}

	response := map[string]interface{}{
		"user_id":      user_info["id"],
		"display_name": user_info["display_name"],
	}

	return response, nil
}

func add_music_to_queue_spotify_func(room_code string, song_name string) (map[string]interface{}, error) {
	user_id, err := get_room_owner_id_db(room_code)
	if err != nil {
		return nil, fmt.Errorf("Failed to get room owner id in db: %w", err)
	}

	acc_token, err := get_access_token(user_id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get access token: %w", err)
	}

	song_info, err := request_song_info(acc_token, song_name)
	if err != nil {
		return nil, fmt.Errorf("Failed to request song uri: %w", err)
	}

	device_id, err := request_available_device_id(acc_token)
	if err != nil {
		return nil, fmt.Errorf("Failed to request available devices: %w", err)
	}

	err = request_add_music_to_queue_spotify(acc_token, device_id, song_info["uri"].(string))
	if err != nil {
		return nil, fmt.Errorf("Failed to request add music to queue spotify: %w", err)
	}

	response := map[string]interface{}{
		"song_name":    song_info["song_name"],
		"song_artists": song_info["artists"],
	}

	return response, nil
}

func create_user_room_func(user_id string) (map[string]interface{}, error) {

	room_code := generate_room_code()
	err := create_user_room_in_db(user_id, room_code)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user room in data base: %w", err)
	}

	response := map[string]interface{}{
		"room_code": room_code,
	}

	return response, nil
}

func current_song_func(room_code string) (map[string]interface{}, error) {

	user_id, err := get_room_owner_id_db(room_code)
	if err != nil {
		return nil, fmt.Errorf("Failed to get room owner id in db: %w", err)
	}

	acc_token, err := get_access_token(user_id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get access token: %w", err)
	}

	current_song_info, err := request_current_song_info(acc_token)
	if err != nil {
		return nil, fmt.Errorf("Failed to request current song %w", err)
	}

	response := map[string]interface{}{
		"song_name":    current_song_info["song_name"],
		"song_artists": current_song_info["artists"],
		"progress_ms":  current_song_info["progress_ms"],
		"duration_ms":  current_song_info["duration_ms"],
	}

	return response, nil
}

func vote_skip_song_func(room_code string) (map[string]interface{}, error) {
	response := make(map[string]interface{})

	room_song, err := get_room_song_in_db(room_code)
	if err != nil {
		return nil, fmt.Errorf("Failed to get room song: %w", err)
	}

	user_id, err := get_room_owner_id_db(room_code)
	if err != nil {
		return nil, fmt.Errorf("Failed to get room owner id in db: %w", err)
	}

	acc_token, err := get_access_token(user_id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get access token: %w", err)
	}

	current_song_info, err := request_current_song_info(acc_token)
	if err != nil {
		return nil, fmt.Errorf("Failed to get the current song: %w", err)
	}
	current_song := current_song_info["song_name"].(string)

	vote_numbers, err := get_skip_votes_in_db(room_code)
	if err != nil {
		return nil, fmt.Errorf("Failed to get number of vote: %w", err)
	}

	users_number, err := get_users_number_in_db(room_code)
	if err != nil {
		return nil, fmt.Errorf("Failed to get users number: %w", err)
	}

	var votes_nd int
	if users_number%2 == 0 {
		votes_nd = users_number / 2
	} else {
		votes_nd = users_number/2 + 1
	}

	if room_song == "" || room_song != current_song {
		update_current_song_in_db(current_song, room_code)
		if err != nil {
			return nil, fmt.Errorf("Failed to update current song in data base: %w", err)
		}

		vote_numbers = 0
	}

	if vote_numbers+1 >= votes_nd || users_number == 2 {
		update_skip_votes_in_db(0, room_code)
		if err != nil {
			return nil, fmt.Errorf("Failed to update votes in database: %w", err)
		}

		err = request_skip_song(acc_token)
		if err != nil {
			return nil, fmt.Errorf("Failed to request skip song: %w", err)
		}

		response["music_skipped"] = true

	} else {
		update_skip_votes_in_db(vote_numbers+1, room_code)
		if err != nil {
			return nil, fmt.Errorf("Failed to update skip vote in database: %w", err)
		}

		response["music_skipped"] = false
		response["missing_votes"] = votes_nd - (vote_numbers + 1)
	}

	response["song"] = current_song

	return response, nil
}
