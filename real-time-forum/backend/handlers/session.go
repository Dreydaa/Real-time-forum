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

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/session" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("Session cookie not found or expired")

		cookie = &http.Cookie{Name: "session", Value: "dummy"}
	}

	foundVal := cookie.Value

	curr, err := database.CurrentUser(config.Path, foundVal)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	cid := strconv.Itoa(curr.ID)

	resp := structure.Resp{
		Msg: cid + "/" + curr.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
}
