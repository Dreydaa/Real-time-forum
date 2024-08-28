package handlers

import (
	"encoding/json"
	"fmt"
	"forum/backend/config"
	"forum/backend/database"
	"forum/backend/structure"
	"net/http"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		param := r.URL.Query().Get("param")
		data := r.URL.Query().Get("data")
		if param == "" || data == "" {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		comments, err := database.FindCommentByParam(config.Path, param, data)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(comments)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)

	case "POST":

		var newComment structure.Comment

		err := json.NewDecoder(r.Body).Decode(&newComment)
		if err != nil {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		fmt.Println(newComment)

		err = database.NewComment(config.Path, newComment)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		var msg = structure.Resp{Msg: "Sent comment"}

		resp, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	default:
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
