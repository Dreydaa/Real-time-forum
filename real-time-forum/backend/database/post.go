package database

import (
	"database/sql"
	"errors"
	"forum/backend/structure"
	"strconv"
	"time"
)

func NewPost(path string, p structure.Post, u structure.User) error {
	db, err := OpenDB(path)
	if err != nil {
		return err
	}

	defer db.Close()

	dt := time.Now().Format("01-02-2006 15:04:05")

	_, err = db.Exec(AddPost, u.ID, p.Category, p.Title, p.Content, dt)
	if err != nil {
		return err
	}
	return nil
}

func ConvertRowToPost(rows *sql.Rows) ([]structure.Post, error) {
	var posts []structure.Post

	for rows.Next() {
		var p structure.Post

		err := rows.Scan(&p.ID, &p.User_id, &p.Category, &p.Content, &p.Date)
		if err != nil {
			break
		}

		posts = append(posts, p)
	}
	return posts, nil
}

func FindAllPosts(path string) ([]structure.Post, error) {
	db, err := OpenDB(path)
	if err != nil {
		return []structure.Post{}, errors.New("failed to open database")
	}
	defer db.Close()

	rows, err := db.Query(GetAllPost)
	if err != nil {
		return []structure.Post{}, errors.New("failed to find posts")
	}

	posts, err := ConvertRowToPost(rows)
	if err != nil {
		return []structure.Post{}, errors.New("failed to convert")
	}

	return posts, nil
}

func FindPostByParam(path, parameter, data string) ([]structure.Post, error) {
	var q *sql.Rows

	db, err := OpenDB(path)
	if err != nil {
		return []structure.Post{}, errors.New("failed to open database")
	}

	defer db.Close()

	switch parameter {
	case "id":
		i, err := strconv.Atoi(data)
		if err != nil {
			return []structure.Post{}, errors.New("id must be an integer")
		}

		q, err = db.Query(GetPostById, i)
		if err != nil {
			return []structure.Post{}, errors.New("could not find id")
		}
	case "user_id":
		q, err = db.Query(GetAllPostByUser, data)
		if err != nil {
			return []structure.Post{}, errors.New("could not find any posts by that user")
		}
	case "category":
		q, err = db.Query(GetAllPostByCategory, data)
		if err != nil {
			return []structure.Post{}, errors.New("could not find any posts with that category")
		}
	default:
		return []structure.Post{}, errors.New("cannot search by that parameter")
	}

	posts, err := ConvertRowToPost(q)
	if err != nil {
		return []structure.Post{}, errors.New("failed to convert")
	}

	return posts, nil
}
