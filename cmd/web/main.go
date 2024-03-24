package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"

	"queuetify.gbatalheiro.pt/internal/models"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	users          *models.UserModel
	rooms          *models.RoomModel
	sessionManager *scs.SessionManager
	CLIENT_ID      string
	CLIENT_SECRET  string
	ADDRESS        string
	CODE_SIZE      int
}

func main() {
	addr := flag.String("addr", "127.0.0.1:4000", "HTTP network address")
	client_id := flag.String("client_id", "", "Spotify client id")
	client_secret := flag.String("client_secret", "", "Spotify client secret")
	code_size := flag.String("code_size", "6", "Size of the coom code")

	flag.Parse()

	if *client_id == "" {
		fmt.Fprintf(os.Stderr, "Error: the -client_id argument is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *client_secret == "" {
		fmt.Fprintf(os.Stderr, "Error: the -client_secret argument is required\n")
		flag.Usage()
		os.Exit(1)
	}

	code_size_num, err := strconv.Atoi(*code_size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: the -code_size argument is invalid must be and integer\n")
		flag.Usage()
		os.Exit(1)
	}

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	sessionManager := scs.New()
	sessionManager.Store = memstore.New()

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		users:          &models.UserModel{DB: new(sync.Map)},
		rooms:          &models.RoomModel{DB: new(sync.Map)},
		sessionManager: sessionManager,
		ADDRESS:        *addr,
		CLIENT_ID:      *client_id,
		CLIENT_SECRET:  *client_secret,
		CODE_SIZE:      code_size_num,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()

	errorLog.Fatal(err)
}
