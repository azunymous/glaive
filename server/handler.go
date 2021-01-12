package server

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	json "github.com/json-iterator/go"
	"github.com/rs/cors"
	"glaive/adapter/asagi"
	"glaive/board"
	"log"
	"net/http"
	"strconv"
	"time"
)

type API struct {
	LoaderRegistry map[string]*asagi.Loader
	URIRegistry    map[string]Board
	Version        string
	DefaultTime    string
}

func (a *API) Handler() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.AllowAll().Handler)

	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/", a.homePageHandler)
	router.Get("/boards", a.overboardHandler)
	router.Get("/{board}/thread/all", a.getAllThreadsHandler)
	router.Get("/{board}/thread", a.getThreadHandler)

	return router
}

func getThreads(loader *asagi.Loader, time time.Time) ([]*board.Thread, error) {
	posts, err := loader.GetPosts(time)
	if err != nil {
		return nil, err
	}
	threads := loader.PostToThreads(posts, asagi.DiscardIfAfter(time))
	return threads, nil
}

func (a *API) homePageHandler(w http.ResponseWriter, _ *http.Request) {
	addHeaders(w)
	_, _ = fmt.Fprintf(w, `{"V" : "%s", "data" : "GLAIVE API"}`, a.Version)
}

func (a *API) overboardHandler(w http.ResponseWriter, _ *http.Request) {
	addHeaders(w)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(a.URIRegistry)
}

func (a *API) getAllThreadsHandler(w http.ResponseWriter, r *http.Request) {
	atUnixSec := r.URL.Query().Get("time")
	if atUnixSec == "" {
		atUnixSec = a.DefaultTime
	}
	log.Printf("Processing get all threads query for %s", atUnixSec)

	unixSecs, err := strconv.Atoi(atUnixSec)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(boardResponse{Status: "FAILURE"})
		return
	}

	boardLoader, ok := a.LoaderRegistry[chi.URLParam(r, "board")]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(boardResponse{Status: "FAILURE"})
		return
	}

	t, err := getThreads(boardLoader, time.Unix(int64(unixSecs), 0))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(boardResponse{Status: "FAILURE"})
		return
	}

	addHeaders(w)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(t)

}

func (a *API) getThreadHandler(w http.ResponseWriter, r *http.Request) {
	threadNo := r.URL.Query().Get("no")
	if threadNo == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(boardResponse{Status: "FAILURE"})
		return
	}

	no, err := strconv.Atoi(threadNo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(boardResponse{Status: "FAILURE"})
		return
	}

	boardLoader, ok := a.LoaderRegistry[chi.URLParam(r, "board")]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(boardResponse{Status: "FAILURE"})
		return
	}
	posts, err := boardLoader.GetPostsByThread(no)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(boardResponse{Status: "FAILURE"})
		return
	}

	t := boardLoader.PostToThread(posts)

	addHeaders(w)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(boardResponse{Status: "SUCCESS", No: threadNo, Thread: *t, Type: "THREAD"})
}

func addHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
