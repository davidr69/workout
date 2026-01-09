package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"workout.lavacro.net/api"
	"workout.lavacro.net/database"
)

type config struct {
	port int
}

func main() {
	fmt.Println("Hello, World!")

	db := &database.Dao{}
	db.Init()

	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "port to listen on")
	flag.Parse()

	mux := api.Routes(db)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/", http.StripPrefix("", fileServer))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	fmt.Println(err)
}
