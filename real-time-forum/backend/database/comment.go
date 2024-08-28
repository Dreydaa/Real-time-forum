package database

import (
	"database/sql"
	"errors"
	"forum/backend/structure"
	"strconv"
	"time"
)

func NewComment(path string, c structure.Comment) error {
	db, err := OpenDB(path)
	if err != nil {
		return err
	}

	defer db.Close()

	dt := time.Now().Format("01-02-2006 15:04:05")

	_, err = db.Exec(AddComment, c.Post_id, c.User_id, c.Content, dt)
	if err != nil {
		return err
	}
	return nil
}

func ConvertRowToComment(rows *sql.Rows) ([]structure.Comment, error) {
	var comments []structure.Comment

	for rows.Next() {
		var c structure.Comment

		err := rows.Scan(&c.ID, &c.Post_id, &c.User_id, &c.Content, &c.Date)
		if err != nil {
			break
		}

		comments = append(comments, c)
	}
	return comments, nil
}

func FindCommentByParam(path, param, data string) ([]structure.Comment, error) {
	var q *sql.Rows

	db, err := OpenDB(path)
	if err != nil {
		return []structure.Comment{}, errors.New("failed to open database")
	}

	defer db.Close()

	i, err := strconv.Atoi(data)
	if err != nil {
		return []structure.Comment{}, errors.New("must provide an integer")
	}

	switch param {
	case "id":

		q, err = db.Query(GetCommentById, i)
		if err != nil {
			return []structure.Comment{}, errors.New("could not find id")
		}
	case "post_id":
		q, err = db.Query(GetAllPostComment, i)
		if err != nil {
			return []structure.Comment{}, errors.New("could not find post_id")
		}
	case "user_id":
		q, err = db.Query(GetAllUserComment, i)
		if err != nil {
			return []structure.Comment{}, errors.New("could not find user_id")
		}
	default:
		return []structure.Comment{}, errors.New("cannot search by that parameter")
	}

	comments, err := ConvertRowToComment(q)
	if err != nil {
		return []structure.Comment{}, errors.New("failed to convert")
	}

	return comments, nil
}
