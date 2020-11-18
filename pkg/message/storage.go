package message

import (
	"database/sql"
)

type Repository interface {
	GetById(id int64) (Chat, error)
	Add(chat Chat) error
	GetAll() ([]Chat, error)
}

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (r *Storage) GetById(id int64) (Chat, error) {
	row := r.db.QueryRow("SELECT * FROM chat WHERE id = ?", id)

	chat := Chat{}

	err := row.Scan(&chat.Id, &chat.Username, &chat.Created)

	return chat, err
}

func (r *Storage) Add(chat Chat) error {
	stmt, err := r.db.Prepare("INSERT INTO chat(id, username) values(?,?)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(chat.Id, chat.Username)

	if err != nil {
		return err
	}

	return nil
}

func (r *Storage) GetAll() ([]Chat, error) {
	chats := make([]Chat, 0)
	rows, err := r.db.Query("SELECT * FROM chat")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		chat := Chat{}

		err := rows.Scan(&chat.Id, &chat.Username, &chat.Created)

		if err != nil {
			return chats, nil
		}

		chats = append(chats, chat)
	}

	return chats, err
}
