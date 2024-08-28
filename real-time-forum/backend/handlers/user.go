package handlers

import (
	"encoding/json"
	"forum/backend/config"
	"forum/backend/database"
	"net/http"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		users, err := database.FindAllUser(config.Path)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	} else {
		user, err := database.FindUserByParam(config.Path, "id", id)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}
