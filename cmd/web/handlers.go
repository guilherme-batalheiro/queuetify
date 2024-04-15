package main

import (
	"html/template"
	"net/http"

	"queuetify.gbatalheiro.pt/internal/spotify"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := struct {
		IsAuthenticated bool
	}{
		IsAuthenticated: app.isAuthenticated(r)}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) userSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	link, err := spotify.GenerateAuthLink(app.CLIENT_ID, app.ADDRESS)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}

func (app *application) userSpotifyLoginCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	code := queryParams.Get("code")
	if code == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	tokens, err := spotify.GetSpotifyAuthTokens(code, app.CLIENT_ID, app.CLIENT_SECRET, app.ADDRESS)
	if err != nil {
		app.errorLog.Println(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	spotifyUserInfo, err := spotify.RequestUserInfo(tokens.AccessToken)
	if err != nil {
		app.errorLog.Println(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	_, exits := app.users.Get(spotifyUserInfo.Id)

	if !exits {
		app.users.Insert(spotifyUserInfo.Id, spotifyUserInfo.DisplayName, spotifyUserInfo.Email)
	}

	ok := app.users.UpdateTokens(spotifyUserInfo.Id, tokens)
	if !ok {
		app.errorLog.Println("Something went wrong in updating tokens!")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", spotifyUserInfo.Id)

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func (app *application) createRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var room_code string
	var ok bool

	spotify_id := app.sessionManager.GetString(r.Context(), "authenticatedUserID")

	room_code, has_room := app.users.GetRoomCode(spotify_id)

	if !has_room {
		room_code = app.generateRoomCode()

		ok = app.users.UpdateRoomCode(spotify_id, room_code)
		if !ok {
			app.errorLog.Println("Something went wrong in adding room code!")
			app.clientError(w, http.StatusBadRequest)
			return
		}

		app.rooms.Insert(room_code, spotify_id)
	}

	http.Redirect(w, r, "/room/"+room_code, http.StatusPermanentRedirect)
}

func (app *application) room(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	room_id := r.PathValue("id")

	room, found := app.rooms.Get(room_id)
	if !found {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/room.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	owner, found := app.users.Get(room.OwnerSpotifyId)
	if !found {
		app.errorLog.Println("Something went wrong is not supposed to have a room withouth owner")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	spotify_id := app.sessionManager.GetString(r.Context(), "authenticatedUserID")

	data := struct {
		RoomCode string
		Owner    string
		IsOwner  bool
	}{
		RoomCode: room_id,
		Owner:    owner.DisplayName,
		IsOwner:  spotify_id == room.OwnerSpotifyId,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) roomDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	room_id := r.PathValue("id")

	room, found := app.rooms.Get(room_id)
	if !found {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	owner, found := app.users.Get(room.OwnerSpotifyId)
	if !found {
		app.errorLog.Println("Something went wrong is not supposed to have a room withouth owner")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	spotify_id := app.sessionManager.GetString(r.Context(), "authenticatedUserID")

	// owner is deleting the room
	if owner.SpotifyId == spotify_id {
		app.rooms.Delete(room_id)
        app.users.UpdateRoomCode(room.OwnerSpotifyId, "")
	}

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func (app *application) roomAddToQueuePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	roomId := r.PathValue("id")

	room, found := app.rooms.Get(roomId)
	if !found {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	owner, found := app.users.Get(room.OwnerSpotifyId)
	if !found {
		app.errorLog.Println("Something went wrong is not supposed to have a room withouth owner")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	songQuery := r.PostForm.Get("song_query")

    accessToken, err := owner.Tokens.GetAccessToken()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	_, err = spotify.RequestAddMusicToQueue(accessToken, songQuery)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}
