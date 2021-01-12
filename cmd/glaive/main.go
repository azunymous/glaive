package main

import (
	"fmt"
	"glaive/adapter/asagi"
	"glaive/server"
	"log"
	"net/http"
	"os"
)

// TODO Replace with a struct and maybe dependency injection
var boardRegistry map[string]*asagi.Loader

func main() {
	fmt.Printf("glaive")
	var err error
	dsn := os.Getenv("IGIARI_SQL_DSN")
	if dsn == "" {
		dsn = "root:mariadbrootpassword@tcp(127.0.0.1:3306)/asagi?charset=utf8&parseTime=True&loc=Local"
	}

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

	srv := server.API{
		BoardRegistry: boardRegistry,
	}

	log.Fatal(http.ListenAndServe(":8080", srv.Handler()))
}
