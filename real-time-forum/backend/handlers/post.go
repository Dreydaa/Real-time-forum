package handlers

import (
	"encoding/json"
	"forum/backend/config"
	"forum/backend/database"
	"forum/backend/structure"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		var posts []structure.Post
		var err error

		param := r.URL.Query().Get("param")
		if param == "" {
			posts, err = database.FindAllPosts(config.Path)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			data := r.URL.Query().Get("data")

			if data == "" {
				http.Error(w, "400 bad request", http.StatusBadRequest)
				return
			}

			posts, err = database.FindPostByParam(config.Path, param, data)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
		}

		resp, err := json.Marshal(posts)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)

	case "POST":

		var newPost structure.Post

		err := json.NewDecoder(r.Body).Decode(&newPost)
		if err != nil {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		cookie, err := r.Cookie("session")
		if err != nil {
			return
		}

		foundVal := cookie.Value

		curr, err := database.CurrentUser(config.Path, foundVal)
		if err != nil {
			return
		}

		err = database.NewPost(config.Path, newPost, curr)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		var msg = structure.Resp{Msg: "New post added"}

		resp, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)

	default:

		http.Error(w, "405 method nor allowed", http.StatusMethodNotAllowed)
		return
	}
}
