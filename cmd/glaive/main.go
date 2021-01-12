package main

import (
	"fmt"
	"glaive/adapter/asagi"
	"glaive/config"
	"glaive/server"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Printf("glaive")
	var err error

	var configPath string
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	conf := config.LoadConfig(configPath)

	db, err := asagi.NewSqlConn(conf.DSN)
	if err != nil {
		log.Fatal(err)
	}

	loaderRegistry := make(map[string]*asagi.Loader, len(conf.Boards))
	uriRegistry := make(map[string]server.Board, len(conf.Boards))
	for name, uris := range conf.Boards {
		letter := strings.TrimSuffix(strings.TrimPrefix(name, "/"), "/")
		loader, err := asagi.NewLoader(letter, db)
		if err != nil {
			log.Fatal(err)
		}
		loaderRegistry[letter] = loader
		uriRegistry[name] = server.Board{
			Host:   uris.URI,
			Images: uris.ImageURI,
		}
	}

	srv := server.API{
		LoaderRegistry: loaderRegistry,
		URIRegistry:    uriRegistry,
		DefaultTime:    conf.DefaultTime,
		Version:        conf.Version,
	}

	log.Fatal(http.ListenAndServe(":8080", srv.Handler()))
}
