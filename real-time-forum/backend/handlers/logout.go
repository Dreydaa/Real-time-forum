package handlers

import (
	"encoding/json"
	"forum/backend/config"
	"forum/backend/database"
	"forum/backend/structure"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := database.OpenDB(config.Path)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	cookie, err := r.Cookie("session")
	if err != nil {
		return
	}

	_, err = db.Exec(database.RemoveCookie, cookie.Value)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	var msg = structure.Resp{Msg: "Ciao"}

	resp, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
