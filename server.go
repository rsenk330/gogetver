package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/stretchr/graceful"
)

// AppConfig holds the configuration for the site.
type AppConfig struct {
	Port              string
	IP                string
	Debug             bool
	Hostname          string
	TemplatesDir      string
	GoogleAnalyticsID string
}

// NewAppConfig creates the configuration from environment variables.
func NewAppConfig() *AppConfig {
	return &AppConfig{
		Hostname:          Getenv("HOSTNAME", "gogetver.com"),
		IP:                Getenv("IP", "127.0.0.1"),
		Port:              Getenv("PORT", "5000"),
		Debug:             GetenvBool("DEBUG", false),
		TemplatesDir:      Getenv("TEMPLATE_ROOT", "./templates"),
		GoogleAnalyticsID: Getenv("GA_TRACKING_ID", ""),
	}
}

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

func Router(app *App) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{pkg:.+}/info/refs", app.GitService).Methods("GET")
	router.HandleFunc("/{pkg:.+}/git-upload-pack", app.GitUploadPack).Methods("POST")
	router.HandleFunc("/{pkg:.+}", app.Package).Methods("GET").Name("package")
	router.HandleFunc("/", app.Home).Methods("GET").Name("home")
	return router
}

func main() {
	app := NewApp(NewAppConfig())

	n := negroni.Classic()
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(Router(app))

	graceful.Run(fmt.Sprintf("%s:%s", app.Config.IP, app.Config.Port), 10*time.Second, n)
}
