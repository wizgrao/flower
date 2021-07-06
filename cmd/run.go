package main

import (
	"github.com/wizgrao/flower"
	"net/http"
)

func main() {
	http.Handle("/ws", flower.NewServer())
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "index.html")
	})
	http.HandleFunc("/flower.js", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "flower.js")
	})
	http.ListenAndServe(":8080", nil)
}