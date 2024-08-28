package main

import (
	"forum/backend/handlers"
	"mime"
	"net/http"
)

func main() {
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".js", "application/javascript")

	static := http.FileServer(http.Dir("frontend"))
	http.Handle("/UI/", http.StripPrefix("/UI/", static))
	http.HandleFunc("/UI/css/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		static.ServeHTTP(w, r)
	})

	http.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		static.ServeHTTP(w, r)
	})

	http.HandleFunc("/chat.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		static.ServeHTTP(w, r)
	})

	handlers.StartServer()
}
