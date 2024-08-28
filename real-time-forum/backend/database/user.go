package database

import (
	"database/sql"
	"errors"
	"strconv"

	"forum/backend/structure"

	_ "github.com/mattn/go-sqlite3"
)

func NewUser(path string, u structure.User) error {
	db, err := OpenDB(path)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(AddUser, u.Username, u.Firstname, u.Lastname, u.Gender, u.Email, u.Password)
	if err != nil {
		return err
	}
	return nil
}

func UserExists(path, value string) (bool, error) {
	db, err := OpenDB(path)
	if err != nil {
		return false, err
	}
	defer db.Close()

	query := `SELECT COUNT(*) FROM users WHERE email = ? OR username = ?`
	row := db.QueryRow(query, value, value)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func ConvertRowToUser(rows *sql.Rows) ([]structure.User, error) {
	var users []structure.User

	for rows.Next() {
		var u structure.User

		err := rows.Scan(&u.ID, &u.Username, &u.Firstname, &u.Lastname, &u.Gender, &u.Email, &u.Password)
		if err != nil {
			break
		}

		users = append(users, u)
	}

	return users, nil
}

func FindAllUser(path string) ([]structure.User, error) {
	db, err := OpenDB(path)
	if err != nil {
		return []structure.User{}, errors.New("failed to open database")
	}

	defer db.Close()

	rows, err := db.Query(GetAllUser)
	if err != nil {
		return []structure.User{}, errors.New("failed to find users")
	}

	users, err := ConvertRowToUser(rows)
	if err != nil {
		return []structure.User{}, errors.New("failed to convert")
	}

	return users, nil
}

func FindUserByParam(path, parameter, data string) (structure.User, error) {
	var q *sql.Rows

	db, err := OpenDB(path)
	if err != nil {
		return structure.User{}, errors.New("failed to open database")
	}

	defer db.Close()

	switch parameter {
	case "id":
		i, err := strconv.Atoi(data)
		if err != nil {
			return structure.User{}, errors.New("id must be an integer")
		}

		q, err = db.Query(GetUserByID, i)
		if err != nil {
			return structure.User{}, errors.New("could not find id")
		}
	case "username":
		q, err = db.Query(GetUserByUsername, data)
		if err != nil {
			return structure.User{}, errors.New("could not find username")
		}
	case "email":
		q, err = db.Query(GetUserByEmail, data)
		if err != nil {
			return structure.User{}, errors.New("could not find email")
		}
	default:
		return structure.User{}, errors.New("cannot search by that parameter")
	}
	user, err := ConvertRowToUser(q)
	if err != nil {
		return structure.User{}, errors.New("failed to convert")
	}
	return user[0], nil
}

func CurrentUser(path, val string) (structure.User, error) {
	db, err := OpenDB(path)
	if err != nil {
		return structure.User{}, errors.New("failed to open database")
	}

	defer db.Close()

	q, err := db.Query(GetSessionUser, val)
	if err != nil {
		return structure.User{}, err
	}

	users, err := ConvertRowToUser(q)
	if err != nil {
		return structure.User{}, err
	}
	if len(users) == 0 {
		return structure.User{}, errors.New("user not found")
	}
	return users[0], nil
}
