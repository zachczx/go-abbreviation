package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go-abbreviation/core"
	"go-abbreviation/models"
	"go-abbreviation/templates"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(middleware.Compress(6))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		core.TemplRender(w, r, templates.MainPage("Home"))
	})
	r.Post("/listlegacy", core.ShowListJson)
	r.Get("/check", func(w http.ResponseWriter, r *http.Request) {
		jsonFile, err := os.Open("static/data.json")
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Println("Successfully opened json file")
		defer jsonFile.Close()

		// Decode json file
		var abv models.List
		json.NewDecoder(jsonFile).Decode(&abv)
		core.TemplRender(w, r, templates.CheckDuplicate("Check Duplicates", abv))
	})
	r.Get("/search", core.ShowListDb)
	r.Get("/all/", core.ShowListDbAlphabets)
	r.Get("/{alphabet}", core.ShowListDbFilterAlphabets)
	r.Get("/list/{alphabet}", func(w http.ResponseWriter, r *http.Request) {
		// Previously used route - /list/{alphabet}, setting redirect to catch those
		param := chi.URLParam(r, "alphabet")
		http.Redirect(w, r, "/"+param, http.StatusMovedPermanently)
	})
	r.Get("/syncjsontodb", core.SyncJsonToDb)
	r.Get("/test", core.Test)

	// Create a route along /static that will serve contents from
	// the ./static/ folder. In lieu of using Nginx to serve.
	// No need for this because I'm using a docker file
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(r, "/static", filesDir)

	listenAddr := os.Getenv("LISTEN_ADDR") //":5173"
	fmt.Println(listenAddr)
	slog.Info("HTTP server started", "listenAddr", listenAddr)
	http.ListenAndServe(listenAddr, r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
