package handlers

import (
	"fmt"
	"net/http"

	"forum/backend/chat"
	"forum/backend/config"
	"forum/backend/database"
)

func StartServer() {
	database.InitDB(config.Path)

	mux := http.NewServeMux()
	hub := chat.NewHub()

	go hub.Run()

	mux.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("/frontend"))))

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/session", SessionHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/user", UserHandler)
	http.HandleFunc("/post", PostHandler)
	http.HandleFunc("/comment", CommentHandler)
	http.HandleFunc("/message", MessageHandler)
	http.HandleFunc("chat", ChatHandler)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, w, r)
	})

	fmt.Println("Server running on port 7001 ✎....")
	http.ListenAndServe(":7001", nil)
	fmt.Println("http://localhost:7001 ✎....")

	if err := http.ListenAndServe(":7001", mux); err != nil {
		fmt.Println("error in main func", err)
	}
}
