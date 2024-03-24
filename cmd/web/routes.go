package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.home)))
	mux.Handle("/userSpotifyLogin/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSpotifyLogin)))
	mux.Handle("/userSpotifyLoginCallback/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSpotifyLoginCallback)))
	mux.Handle("/room/{id}/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.room)))
	mux.Handle("/room/createRoom", app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.createRoom))))
	// mux.HandleFunc("/room/add_to_queue", app.snippetView)
	// mux.HandleFunc("/room/join_room", app.snippetCreate)
	// mux.HandleFunc("/room/exit_room", app.snippetCreate)
	// mux.HandleFunc("/room/vote_skip_song", app.snippetCreate)
	// mux.HandleFunc("/room/current_song", app.snippetCreate)
	// mux.HandleFunc("/room/delete_room", app.snippetCreate)

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
