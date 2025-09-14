package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)


func main(){
	r := chi.NewRouter()


	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world")
	})


	addr := ":2000"
	log.Printf("Server listening on %s\n", addr)

	err := http.ListenAndServe(addr, r)
	if  err != nil {
		log.Fatal(err)
	}

}