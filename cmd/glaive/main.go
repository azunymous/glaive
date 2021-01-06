package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"glaive/adapter/asagi"
	"glaive/board"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var loader *asagi.Loader
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	fmt.Printf("glaive")
	var err error
	loader, err = asagi.NewLoader("root:mariadbrootpassword@tcp(127.0.0.1:3306)/asagi?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", handler()))
}

func getThreads(loader *asagi.Loader, time time.Time) ([]*board.Thread, error) {
	posts, err := loader.GetPosts(time)
	if err != nil {
		return nil, err
	}
	threads := asagi.PostToThreads(posts, asagi.DiscardIfAfter(time))
	return threads, nil
}

func handler() http.Handler {
	router := httprouter.New()
	router.GET("/", homePageHandler)
	router.GET("/boards", overboardHandler)
	router.GET("/thread/all", getAllThreadsHandler)
	router.GET("/thread", getThreadHandler)

	return cors.AllowAll().Handler(router)
}

func homePageHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	addHeaders(w)
	_, _ = fmt.Fprintf(w, `{"V" : "1", "data" : "GLAIVE API"}`)
}

func overboardHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	addHeaders(w)
	var boards = map[string]Board{
		"/c/": {Host: "http://localhost:8080", Images: "/img/"},
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(boards)
}

func getAllThreadsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	t, err := getThreads(loader, time.Unix(int64(unixSecs), 0))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(boardResponse{Status: "FAILURE"})
		return
	}

	addHeaders(w)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(t)

}

func getThreadHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	posts, err := loader.GetPostsByThread(no)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(boardResponse{Status: "FAILURE"})
		return
	}

	t := asagi.PostToThread(posts)

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
