package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Sigafoos/pokewants/gamemaster"
	"github.com/Sigafoos/pokewants/handler"
	"github.com/Sigafoos/pokewants/logger"
	"github.com/Sigafoos/pokewants/middleware"
	"github.com/Sigafoos/pokewants/wants"

	"github.com/NYTimes/gziphandler"
	"github.com/gocraft/dbr"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("no dsn found")
	}

	db, err := dbr.Open("postgres", dsn, nil)
	if err != nil {
		log.Panicln("cannot open database: " + err.Error())
	}
	defer db.Close()

	gm := gamemaster.New(&http.Client{})

	w, err := wants.New(db, logger.New(os.Stderr), gm)
	if err != nil {
		log.Println(err)
		return
	}

	h := handler.New(w)
	mux := http.NewServeMux()
	mux.Handle("/want", http.HandlerFunc(h.HandleWant))
	mux.Handle("/search", http.HandlerFunc(h.HandleSearch))

	chain := gziphandler.GzipHandler(mux)
	chain = middleware.UseJSON(chain)
	chain = middleware.UseAuth(chain)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}
	fmt.Println("server running on port " + port)
	fmt.Println(server.ListenAndServe())
}
