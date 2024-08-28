package handlers

import (
	"encoding/json"
	"forum/backend/config"
	"forum/backend/database"
	"forum/backend/structure"
	"net/http"
	"net/mail"
	"strconv"

	uuid "github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
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

	var loginData structure.Login

	err = json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}

	var param string

	if _, err := mail.ParseAddress(loginData.Data); err != nil {
		param = "username"
	} else {
		param = "email"
	}

	foundUser, err := database.FindUserByParam(config.Path, param, loginData.Data)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginData.Password))
	if err != nil {
		http.Error(w, "401 unauthorized: username or password incorrect", http.StatusUnauthorized)
		return
	}

	_, err = db.Exec(database.RemoveCookie, foundUser.ID)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		sessionId, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		cookie = &http.Cookie{
			Name:     "session",
			Value:    sessionId.String(),
			HttpOnly: true,
			Path:     "/",
			MaxAge:   config.CookieAge,
			SameSite: http.SameSiteDefaultMode,
			Secure:   true,
		}
		http.SetCookie(w, cookie)
	}

	_, err = db.Exec(database.AddSessions, cookie.Value, foundUser.ID)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	cid := strconv.Itoa(foundUser.ID)

	var msg = structure.Resp{Msg: cid + "/" + foundUser.Username}

	resp, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
