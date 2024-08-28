package database

import (
	"database/sql"
	"errors"
	"forum/backend/structure"
	"strconv"
)

func NewMessage(path string, m structure.Message) error {
	db, err := OpenDB(path)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(AddMessage, m.Sender_id, m.Receiver_id, m.Content, m.Date)
	if err != nil {
		return err
	}

	err = UpdateChatTime(m.Sender_id, m.Receiver_id, db)
	if err != nil {
		return err
	}

	return nil
}

func ConvertRowToMessage(rows *sql.Rows) ([]structure.Message, error) {
	var messages []structure.Message

	for rows.Next() {
		var m structure.Message

		err := rows.Scan(&m.ID, &m.Sender_id, &m.Receiver_id, &m.Content, &m.Date)
		if err != nil {
			break
		}

		messages = append(messages, m)
	}

	return messages, nil
}

func FindChatMessages(path, sender, receiver string, firstId int) ([]structure.Message, error) {
	db, err := OpenDB(path)
	if err != nil {
		return []structure.Message{}, errors.New("failed to open database")
	}

	defer db.Close()

	// coversions sender et receiver IDs en integers
	s, err := strconv.Atoi(sender)
	if err != nil {
		return []structure.Message{}, errors.New("sender id must be an integer")
	}

	r, err := strconv.Atoi(receiver)
	if err != nil {
		return []structure.Message{}, errors.New("receiver id must be an integer")
	}

	q, err := db.Query(GetAllChatMessage, s, r, r, s, firstId)

	//`SELECT * FROM messages WHERE ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND ( id <= ? ) ORDER BY id DESC LIMIT 10`
	if err != nil {
		return []structure.Message{}, errors.New("could not find chat messages")
	}

	messages, err := ConvertRowToMessage(q)
	if err != nil {
		return []structure.Message{}, errors.New("failed to convert")
	}

	return messages, nil
}

func FindLastMessage(path, sender, receiver string) (structure.Message, error) {
	db, err := OpenDB(path)
	if err != nil {
		return structure.Message{}, errors.New("failed to open database")
	}

	defer db.Close()

	s, err := strconv.Atoi(sender)
	if err != nil {
		return structure.Message{}, errors.New("sender id must be an integer")
	}

	r, err := strconv.Atoi(receiver)
	if err != nil {
		return structure.Message{}, errors.New("receiver id must be an integer")
	}

	q, err := db.Query(GetLastMessage, r, r, r, s)

	//`SELECT * FROM messages WHERE ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) ORDER BY id DESC LIMIT 1`
	if err != nil {
		return structure.Message{}, errors.New("could not find chat messages")
	}

	messages, err := ConvertRowToMessage(q)
	if err != nil {
		return structure.Message{}, errors.New("failed to convert")
	}

	return messages[0], nil
}
