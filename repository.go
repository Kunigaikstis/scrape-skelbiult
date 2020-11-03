package main

import (
	"database/sql"
	"time"
)

type Listing struct {
	Id      int
	Url     string
	Created time.Time
}

type ListingRepository struct {
	db *sql.DB
}

type Chat struct {
	Id       int64
	Username string
	Created  time.Time
}

type ChatRepository struct {
	db *sql.DB
}

func NewListingRepository(db *sql.DB) *ListingRepository {
	return &ListingRepository{db: db}
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ListingRepository) GetById(id int) (Listing, error) {
	row := r.db.QueryRow("SELECT * FROM listing WHERE id = ?", id)

	listing := Listing{}

	err := row.Scan(&listing.Id, &listing.Url, &listing.Created)

	return listing, err
}

func (r *ListingRepository) Save(listing Listing) error {
	stmt, err := r.db.Prepare("INSERT INTO listing(id, url) values(?,?)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(listing.Id, listing.Url)

	if err != nil {
		return err
	}

	return nil
}

func (r *ChatRepository) GetById(id int64) (Chat, error) {
	row := r.db.QueryRow("SELECT * FROM chat WHERE id = ?", id)

	chat := Chat{}

	err := row.Scan(&chat.Id, &chat.Username, &chat.Created)

	return chat, err
}

func (r *ChatRepository) Save(chat Chat) error {
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

func (r *ChatRepository) GetAll() ([]Chat, error) {
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
