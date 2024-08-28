package handlers

import (
	"encoding/json"
	"forum/backend/config"
	"forum/backend/database"
	"forum/backend/structure"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUser structure.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "400 Bad request", http.StatusBadRequest)
		return
	}

	if !isValidEmail(newUser.Email) {
		http.Error(w, "400 Bad request : Invalid email address", http.StatusBadRequest)
		return
	}

	emailExist, err := database.UserExists(config.Path, newUser.Email)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	usernameExists, err := database.UserExists(config.Path, newUser.Username)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	if emailExist && usernameExists {
		http.Error(w, "400 Bad request : Email or username already exists", http.StatusBadRequest)
		return
	} else if emailExist {
		http.Error(w, "409 conflict : the email is already taken", http.StatusConflict)
		return
	} else if usernameExists {
		http.Error(w, "409 conflict : the username is already taken", http.StatusConflict)
		return
	}

	passwordHash, err := GenerateHash(newUser.Password)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	newUser.Password = passwordHash

	err = database.NewUser(config.Path, newUser)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	var msg = structure.Resp{Msg: "Successful registration"}

	resp, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)

	return string(hash), err
}

func isValidEmail(email string) bool {
	// Email validation pattern using regular expression
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}
