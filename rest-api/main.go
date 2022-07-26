package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func add_music_to_queue_spotify(c *gin.Context) {
	// add music to queue in spotify
	c.Header("Access-Control-Allow-Origin", os.Getenv("CORS"))
	room_code := c.Request.URL.Query().Get("room_code")
	song_name := c.Request.URL.Query().Get("song_name")
	song_name = strings.ReplaceAll(song_name, "%20", " ")

	response, err := add_music_to_queue_spotify_func(room_code, song_name)
	if err != nil {
		log.Println("Error: Failed to add music to queue of spotify:", err)
		c.IndentedJSON(http.StatusBadGateway, "failed to add music to queue of spotify!")
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func delete_user_room(c *gin.Context) {
	// delete user room
	c.Header("Access-Control-Allow-Origin", os.Getenv("CORS"))
	user_id := c.Request.URL.Query().Get("user_id")

	err := delete_user_room_db(user_id)
	if err != nil {
		log.Println("Error: Ffailed to delete room in database:", err)
		c.IndentedJSON(http.StatusBadGateway, "failed to delete user room!")
		return
	}
}

func create_user_room(c *gin.Context) {
	// create user room
	c.Header("Access-Control-Allow-Origin", os.Getenv("CORS"))
	user_id := c.Request.URL.Query().Get("user_id")

	response, err := create_user_room_func(user_id)
	if err != nil {
		log.Println("Error: Failed to  create user room:", err)
		c.IndentedJSON(http.StatusBadGateway, "failed to create user room!")
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func join_room(c *gin.Context) {
	// join a room
	c.Header("Access-Control-Allow-Origin", os.Getenv("CORS"))
	room_code := c.Request.URL.Query().Get("room_code")

	err := icrement_users_number_in_db(room_code)
	if err != nil {
		log.Println("Error: Failed to join room:", err)
		c.IndentedJSON(http.StatusBadGateway, "failed to join room!")
		return
	}

}

func exit_room(c *gin.Context) {
	// exit room
	c.Header("Access-Control-Allow-Origin", os.Getenv("CORS"))
	room_code := c.Request.URL.Query().Get("room_code")

	err := decrement_users_number_in_db(room_code)
	if err != nil {
		log.Println("Error: Failed to exit room:", err)
		c.IndentedJSON(http.StatusBadGateway, "failed to exit room!")
		return
	}

}

func auth(c *gin.Context) {
	// log in into the app with spotify
	c.Header("Access-Control-Allow-Origin", os.Getenv("CORS"))
	code := c.Request.URL.Query().Get("code")

	response, err := auth_func(code)
	if err != nil {
		log.Println("Error: Failed to log in:", err)
		c.IndentedJSON(http.StatusBadGateway, "failed to log in!")
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func current_song(c *gin.Context) {
	// return current song playing
	c.Header("Access-Control-Allow-Origin", os.Getenv("CORS"))
	room_code := c.Request.URL.Query().Get("room_code")

	response, err := current_song_func(room_code)
	if err != nil {
		log.Println("Error: Failed to get the current song:", err)
		c.IndentedJSON(http.StatusBadGateway, "failed to get the current song!")
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func vote_skip_song(c *gin.Context) {
	// vote to skip song
	c.Header("Access-Control-Allow-Origin", os.Getenv("CORS"))
	room_code := c.Request.URL.Query().Get("room_code")

	response, err := vote_skip_song_func(room_code)
	if err != nil {
		log.Println("Error: Failed to vote to skip song:", err)
		c.IndentedJSON(http.StatusBadGateway, "failed to vote!")
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func checkAPI(c *gin.Context) {
	// sends a ok message to the client if the api is working

	c.IndentedJSON(http.StatusOK, "ok")
}

func checkDB(c *gin.Context) {
	// check if database is available exit api if it fails

	err := db.Ping()
	if err != nil {
		defer log.Fatal("Error: Failed to ping data base")
		c.IndentedJSON(http.StatusBadGateway, "failed to pind data base")
	}
	c.IndentedJSON(http.StatusOK, "ok")
}

func main() {
	var err error

	// start database
	db, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// create users table in database
	err = create_users_table()
	if err != nil && err.Error() != "table users already exists" {
		log.Fatal(err)
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// api endpoints
	router := gin.Default()
	router.GET("/auth", auth)
	router.GET("/add_to_queue", add_music_to_queue_spotify)
	router.POST("/join_room", join_room)
	router.POST("/exit_room", exit_room)
	router.GET("/create_room", create_user_room)
	router.POST("/delete_room", delete_user_room)
	router.GET("/vote_skip_song", vote_skip_song)
	router.GET("/current_song", current_song)
	router.GET("/pingAPI", checkAPI)
	router.GET("/pingDB", checkDB)
	router.Run(os.Getenv("ADDRESS"))
}
