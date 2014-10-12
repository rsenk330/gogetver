package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/stretchr/graceful"
)

// Getenv is just a simple helper to get an environment variable,
// or the default if it is not set.
func Getenv(name, def string) string {
	val := os.Getenv(name)
	if val == "" {
		val = def
	}
	return val
}

// GetenvBool gets an environment variable that is epected to be a
// boolean value.
func GetenvBool(name string, def bool) bool {
	val := Getenv(name, "")
	if val == "" {
		return def
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		boolVal = false
	}
	return boolVal
}

func main() {
	app := NewApp()
	router := mux.NewRouter()
	router.HandleFunc("/{pkg:.+}/info/refs", app.GitService).Methods("GET")
	router.HandleFunc("/{pkg:.+}/git-upload-pack", app.GitUploadPack).Methods("POST")
	router.HandleFunc("/{pkg:.+}", app.Package).Methods("GET")

	port := Getenv("PORT", "5000")
	ip := Getenv("IP", "127.0.0.1")

	n := negroni.Classic()
	n.UseHandler(router)

	graceful.Run(fmt.Sprintf("%s:%s", ip, port), 10*time.Second, n)
}
