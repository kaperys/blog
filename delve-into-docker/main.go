package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		message := r.URL.Path
		message = strings.TrimPrefix(message, "/")
		message = "Hello, " + message + "!"

		w.Write([]byte(message))
	})

	log.Print("starting web server")
	if err := http.ListenAndServe(":1541", nil); err != nil {
		log.Fatal(err)
	}
}
