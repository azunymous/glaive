package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	jsoniter "github.com/json-iterator/go"
	"glaive/adapter/asagi"
	"glaive/board"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/cors"
)

// TODO Replace with a struct and maybe dependency injection
var boardRegistry map[string]*asagi.Loader
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	fmt.Printf("glaive")
	var err error

	const dsn = "root:mariadbrootpassword@tcp(127.0.0.1:3306)/asagi?charset=utf8&parseTime=True&loc=Local"

	db, err := asagi.NewSqlConn(dsn)
	if err != nil {
		log.Fatal(err)
	}
	cloader, err := asagi.NewLoader("c", db)
	if err != nil {
		log.Fatal(err)
	}
	poloader, err := asagi.NewLoader("po", db)
	if err != nil {
		log.Fatal(err)
	}

	boardRegistry = map[string]*asagi.Loader{
		"c":  cloader,
		"po": poloader,
	}

	log.Fatal(http.ListenAndServe(":8080", handler()))
}

func getThreads(loader *asagi.Loader, time time.Time) ([]*board.Thread, error) {
	posts, err := loader.GetPosts(time)
	if err != nil {
		return nil, err
	}
	threads := loader.PostToThreads(posts, asagi.DiscardIfAfter(time))
	return threads, nil
}

func handler() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.AllowAll().Handler)

	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/", homePageHandler)
	router.Get("/boards", overboardHandler)
	router.Get("/{board}/thread/all", getAllThreadsHandler)
	router.Get("/{board}/thread", getThreadHandler)

	return router
}

func homePageHandler(w http.ResponseWriter, _ *http.Request) {
	addHeaders(w)
	_, _ = fmt.Fprintf(w, `{"V" : "1", "data" : "GLAIVE API"}`)
}

func overboardHandler(w http.ResponseWriter, _ *http.Request) {
	addHeaders(w)
	var boards = map[string]Board{
		"/c/":  {Host: "http://localhost:8080/c", Images: "/img/"},
		"/po/": {Host: "http://localhost:8080/po", Images: "/img/"},
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(boards)
}

func getAllThreadsHandler(w http.ResponseWriter, r *http.Request) {
	atUnixSec := r.URL.Query().Get("time")
	if atUnixSec == "" {
		atUnixSec = "1343080585"
	}
	log.Printf("Processing get all threads query for %s", atUnixSec)

	unixSecs, err := strconv.Atoi(atUnixSec)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(boardResponse{Status: "FAILURE"})
		return
	}

	boardLoader, ok := boardRegistry[chi.URLParam(r, "board")]
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

func getThreadHandler(w http.ResponseWriter, r *http.Request) {
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

	boardLoader, ok := boardRegistry[chi.URLParam(r, "board")]
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

type boardResponse struct {
	Status string       `json:"status"`
	No     string       `json:"no"`
	Thread board.Thread `json:"thread"`
	Type   string       `json:"type"`
}

type Board struct {
	Host   string `json:"host"`
	Images string `json:"images"`
}
