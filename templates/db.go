package main

import (
	"database/sql"
	"fmt"
	"time"
)

var db *sql.DB

func get_users_number_in_db(room_code string) (int, error) {
	var users_number int

	err := db.QueryRow("SELECT users_number FROM users WHERE room_code = ?", room_code).Scan(&users_number)
	if err != nil {
		return -1, err
	}

	return users_number, nil
}

func get_room_song_in_db(room_code string) (string, error) {
	var current_song sql.NullString

	err := db.QueryRow("SELECT current_song FROM users WHERE room_code = ?", room_code).Scan(&current_song)
	if err != nil {
		return "", err
	}

	return current_song.String, nil
}

func get_skip_votes_in_db(room_code string) (int, error) {
	var skip_votes int

	err := db.QueryRow("SELECT skip_votes From users WHERE room_code = ?", room_code).Scan(&skip_votes)
	if err != nil {
		return -1, err
	}

	return skip_votes, nil
}

func update_skip_votes_in_db(votes_number int, room_code string) error {
	stmt, err := db.Prepare("UPDATE users SET skip_votes = ? WHERE room_code = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(votes_number, room_code)
	if err != nil {
		return err
	}

	return nil
}

func update_current_song_in_db(current_song string, room_code string) error {
	stmt, err := db.Prepare("UPDATE users SET current_song = ? WHERE room_code = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(current_song, room_code)
	if err != nil {
		return err
	}

	return nil
}

func icrement_users_number_in_db(room_code string) error {
	stmt, err := db.Prepare("UPDATE users SET users_number = users_number + 1 WHERE room_code = ?")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(room_code)
	if err != nil {
		return err
	}

	rows_affected, err := result.RowsAffected()
	// check if a row was affected this is used to check if the room_code is valid
	if err != nil {
		return err
	}

	if rows_affected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func decrement_users_number_in_db(room_code string) error {
	stmt, err := db.Prepare("UPDATE users SET users_number = users_number - 1 WHERE room_code = ?")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(room_code)
	if err != nil {
		return err
	}

	rows_affected, err := result.RowsAffected()
	// check if a row was affected this is used to check if the room_code is valid
	if err != nil {
		return err
	}

	if rows_affected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func get_room_owner_id_db(room_code string) (string, error) {
	var owner_id string

	err := db.QueryRow("SELECT id FROM users WHERE room_code = ?", room_code).Scan(&owner_id)
	if err != nil {
		return "", err
	}

	return owner_id, nil
}

func delete_user_room_db(user_id string) error {
	stmt, err := db.Prepare("UPDATE users SET room_code = NULL, users_number = 0, skip_votes = 0 WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user_id)
	if err != nil {
		return err
	}

	return nil
}

func create_user_room_in_db(user_id string, room_code string) error {

	stmt, err := db.Prepare("UPDATE users SET room_code = ?, users_number = 1, skip_votes = 0 WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(room_code, user_id)
	if err != nil {
		return err
	}

	return nil
}

func update_access_token_in_db(user_id string, token map[string]interface{}) error {
	stmt, err := db.Prepare("UPDATE users SET access_token = ?, access_token_exp = ? WHERE id = ?")
	if err != nil {
		return err
	}

	current_time := time.Now().Unix()
	expires_date := current_time + int64(token["expires_in"].(float64))

	_, err = stmt.Exec(token["access_token"], expires_date, user_id)
	if err != nil {
		return err
	}

	return nil
}

func get_refresh_token_from_db(user_id string) (string, error) {
	var refre_token string

	err := db.QueryRow("SELECT refresh_token FROM users WHERE id = ?", user_id).Scan(&refre_token)
	if err != nil {
		return "", err
	}

	return refre_token, nil
}

func get_access_token_from_db(user_id string) (string, error) {
	var acc_token string

	err := db.QueryRow("SELECT access_token FROM users WHERE id = ?", user_id).Scan(&acc_token)
	if err != nil {
		return "", err
	}

	return acc_token, nil
}

func check_if_user_exits_in_db(user_id string) (bool, error) {
	var id string

	err := db.QueryRow("SELECT id FROM users WHERE id = ?", user_id).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func access_token_is_valid_db(user_id string) (bool, error) {
	var valid bool
	current_t := time.Now().Unix()

	err := db.QueryRow("SELECT (access_token_exp > ?) FROM users WHERE id = ?", current_t, user_id).Scan(&valid)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}

		return false, nil
	}

	return valid, nil
}

func add_user_to_db(tokens map[string]interface{}, user_info map[string]interface{}) error {
	stmt, err := db.Prepare("INSERT INTO users(id, display_name, email, access_token, access_token_exp, refresh_token) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	current_t := time.Now().Unix()
	expires_date := current_t + int64(tokens["expires_in"].(float64))

	_, err = stmt.Exec(user_info["id"], user_info["display_name"], user_info["email"], tokens["access_token"], expires_date, tokens["refresh_token"])
	if err != nil {
		return err
	}

	return nil
}

func update_user_in_db(tokens map[string]interface{}, user_info map[string]interface{}) error {
	stmt, err := db.Prepare("UPDATE users SET display_name = ?, email = ?, access_token = ?, access_token_exp = ?, refresh_token = ? WHERE id = ?")
	if err != nil {
		return err
	}

	current_t := time.Now().Unix()
	expires_date := current_t + int64(tokens["expires_in"].(float64))

	_, err = stmt.Exec(user_info["display_name"], user_info["email"], tokens["access_token"], expires_date, tokens["refresh_token"], user_info["id"])
	if err != nil {
		return err
	}

	return nil
}

func create_users_table() error {
	// create user table in database
	sqlStmt := `
	create table users (
		id text NOT NULL primary key,
		display_name text,
		email text,
		access_token text,
		access_token_exp int,
		refresh_token text,
		room_code text,
        users_number int,
        current_song text,
        skip_votes int
	);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}
