package handlers

import (
	"encoding/json"
	"fmt"
	"forum/backend/config"
	"forum/backend/database"
	"forum/backend/structure"
	"net/http"
	"strconv"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/message" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		cookie, err := r.Cookie("session")
		if err != nil {
			return
		}

		foundval := cookie.Value

		curr, err := database.CurrentUser(config.Path, foundval)
		if err != nil {
			return
		}

		s := strconv.Itoa(curr.ID)

		firstId, _ := strconv.Atoi(r.URL.Query().Get("firstId"))

		fmt.Println("id", firstId)
		r := r.URL.Query().Get("receiver")

		if r == "" {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		messages, err := database.FindChatMessages(config.Path, s, r, firstId)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(messages)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)

	case "POST":
		var newMessage structure.Message

		err := json.NewDecoder(r.Body).Decode(&newMessage)
		if err != nil {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		err = database.NewMessage(config.Path, newMessage)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}
}
