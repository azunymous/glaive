package main

import (
	"encoding/json"
	"fmt"
	"glaive/adapter/asagi"
	"glaive/board"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var loader *asagi.Loader

func main() {
	fmt.Printf("glaive")
	var err error
	loader, err = asagi.NewLoader("root:mariadbrootpassword@tcp(127.0.0.1:3306)/asagi?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	threads, _, _ := getThreads(loader)
	fmt.Printf("%v", threads[0])

	log.Fatal(http.ListenAndServe(":8080", handler()))
}

func getThreads(loader *asagi.Loader) ([]*board.Thread, map[uint64]*board.Thread, error) {
	posts, err := loader.GetPosts(1343080585)
	if err != nil {
		return nil, nil, err
	}
	threads, m := asagi.PostToThreads(posts)
	return threads, m, nil
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
	_, _ = fmt.Fprintf(w, `{"V" : "1", "data" : "ALICE API"}`)
}

func overboardHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	addHeaders(w)
	var boards = map[string]Board{
		"/c/": {Host: "http://localhost:8080", Images: "/img/"},
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(boards)
}

func getAllThreadsHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	t, _, err := getThreads(loader)

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
	_, m, err := getThreads(loader)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(boardResponse{Status: "FAILURE"})
		return
	}

	no, err := strconv.Atoi(threadNo)
	t, ok := m[uint64(no)]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(boardResponse{Status: "FAILURE"})
		return
	}

	addHeaders(w)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(boardResponse{Status: "SUCCESS", No: threadNo, Thread: *t, Type: "THREAD"})

}

func addHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

type userResponse struct {
	Status   string `json:"status"`
	Username string `json:"username"`
	Error    string `json:"error"`
	Token    string `json:"token"`
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
