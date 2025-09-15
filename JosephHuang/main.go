package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
func main() {
	r := chi.NewRouter() // create a chi router
	r.Use(middleware.Logger) // logs requests and responses
	r.Get("/", rootHandler) // declare route and handler
	r.Post("/upload", McapUpload)

	http.ListenAndServe(":2000", r) // start up server at :3000
	}
	func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

 
func McapUpload (w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Uplaoding"))
}