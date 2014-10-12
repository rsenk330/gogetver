package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/stretchr/graceful"
)

func main() {
	app := NewApp()
	router := mux.NewRouter()
	router.HandleFunc("/{pkg:.+}/info/refs", app.GitService).Methods("GET")
	router.HandleFunc("/{pkg:.+}/git-upload-pack", app.GitUploadPack).Methods("POST")
	router.HandleFunc("/{pkg:.+}", app.Package).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	ip := os.Getenv("IP")
	if ip == "" {
		ip = "127.0.0.1"
	}

	n := negroni.Classic()
	n.UseHandler(router)

	graceful.Run(fmt.Sprintf("%s:%s", ip, port), 10*time.Second, n)
}
