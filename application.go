package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/flosch/pongo2.v2"
)

// ByLength sorts an array of strings by length
type ByLength []string

func (b ByLength) Len() int {
	return len(b)
}
func (b ByLength) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByLength) Less(i, j int) bool {
	return len(b[i]) < len(b[j])
}

// PossibleVersions generates a set of possible version from the url.
func PossibleVersions(url string) []string {
	var versions []string
	for i := 3; i <= strings.Count(url, ".")+1; i++ {
		version := strings.SplitAfterN(url, ".", i)
		versions = append(versions, version[len(version)-1])
	}

	// Try to match the longest strings first
	sort.Sort(sort.Reverse(ByLength(versions)))

	if len(versions) == 0 {
		versions = append(versions, "master")
	}

	return versions
}

// App holds the common items for the app.
type App struct {
	Config   *AppConfig
	Renderer *pongo2.TemplateSet
}

// NewApp creates a new App instance with the default configuration.
func NewApp(config *AppConfig) *App {
	renderer := pongo2.NewSet("templates")

	// Configure the pongo renderer
	renderer.Debug = config.Debug
	renderer.SetBaseDirectory(config.TemplatesDir)
	renderer.Globals["Config"] = config

	return &App{
		Renderer: renderer,
		Config:   config,
	}
}

// Home renders the home page.
func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := app.Renderer.FromCache("home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteWriter(nil, w)
}

// Package is the endpoint for the package detail page. If `?go-get=1` is passed,
// the page with the properly formatted ``<meta name="go-import">`` tag is rendered
// for the `go get` tool.
func (app *App) Package(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["pkg"]

	if r.FormValue("go-get") == "1" {
		tmpl, err := app.Renderer.FromCache("go_import.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		ctx := pongo2.Context{
			"Path": path,
			"Vcs":  "git",
		}

		tmpl.ExecuteWriter(ctx, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	tmpl, err := app.Renderer.FromCache("package.html")
	ctx := pongo2.Context{
		"Path": path,
		"Vcs":  "git",
	}

	tmpl.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) getRefs(path string) ([]byte, error) {
	refsPath := fmt.Sprintf("https://%v.git/info/refs?service=git-upload-pack", path)
	resp, err := http.Get(refsPath)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return nil, errors.New("Repo does not exist.")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (app *App) getVersionRef(path, version string) ([]byte, string, string, error) {
	if version == "" {
		versions := PossibleVersions(path)
		for _, v := range versions {
			pathWithoutVer := strings.Split(path, fmt.Sprintf(".%v", v))[0]
			versionRefs := []string{
				fmt.Sprintf("refs/heads/%v", v),
				fmt.Sprintf("refs/tags/%v", v),
			}

			resp, err := app.getRefs(pathWithoutVer)
			if err == nil {
				strResp := string(resp)

				for _, versionRef := range versionRefs {
					if strings.Contains(strResp, versionRef) {
						return resp, versionRef, pathWithoutVer, nil
					}
				}
			}
		}
	}

	return nil, "", "", errors.New("Could not find version.")
}

// GitService is the endpoint for getting the git refs. It hacks the results
// to trick the `go get` tool into downloading the specified version.
//
// The only allowed service in `git-upload-pack`. Any other service will result in
// a 404.
//
// Documentation on how git handles this can be found here:
// https://github.com/git/git/blob/master/Documentation/technical/http-protocol.txt
func (app *App) GitService(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("service") != "git-upload-pack" {
		http.NotFound(w, r)
		return
	}

	vars := mux.Vars(r)
	path := vars["pkg"]

	refs, matchedVer, _, err := app.getVersionRef(path, "")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	refsStr := string(refs)

	// Get the sha of the version to use
	versionRegexp := regexp.MustCompile(fmt.Sprintf("(?m)(?P<start>^[0-9a-f]{4})([0-9a-f]{40}) %v$", matchedVer))
	versionHash := versionRegexp.FindStringSubmatch(refsStr)[2]

	// Replace the master branch sha with the version to use
	headRegexp := regexp.MustCompile("(?m)(?P<start>[0-9a-f]{4})(?P<hash>[0-9a-f]{40}) refs/heads/master$")
	refsStr = headRegexp.ReplaceAllString(refsStr, fmt.Sprintf("${start}%v refs/heads/master", versionHash))

	w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(refsStr))
}

// GitUploadPack issues a 301 redirect to the actual server to download the package.
func (app *App) GitUploadPack(w http.ResponseWriter, r *http.Request) {
	// http://git-scm.com/book/en/Git-Internals-Git-References
	vars := mux.Vars(r)
	path := vars["pkg"]

	_, _, pathWithoutVer, err := app.getVersionRef(path, "")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("https://%v.git/git-upload-pack", pathWithoutVer), http.StatusMovedPermanently)
}
