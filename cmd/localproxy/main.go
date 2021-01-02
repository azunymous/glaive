package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const port = "8000"

const srv = "http://localhost:3000"
const fileroot = "/mnt/hgfs/progdev/archive/"

func main() {

	u, _ := url.Parse(srv)
	// start server
	http.Handle("/", &baseHandle{
		srv: u,
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type baseHandle struct {
	srv *url.URL
}

func (h *baseHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/img") {
		http.StripPrefix("/img/", http.FileServer(http.Dir(fileroot))).ServeHTTP(w, r)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(h.srv)
	proxy.ServeHTTP(w, r)

}
