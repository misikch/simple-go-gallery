package main

import (
	"database/sql"
	"fmt"
)

type Photo struct {
	ID int
	UserID int
	Path string
}

type StDb struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *StDb {
	return &StDb{
		db: db,
	}
}

func (st *StDb) Add(p *Photo) error {
	//language=MySQL
	res, err := st.db.Exec("INSERT INTO `photos`(`user_id`, `path`) VALUES (?, ?)", p.UserID, p.Path)
	if err != nil {
		return fmt.Errorf("failed to query insert into photos table: %w", err)
	}

	li, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last_insert_id: %w", err)
	}

	if li == 0 {
		return fmt.Errorf("last_insert_id is 0")
	}

	return nil
}

func (st *StDb) GetPhotos(userID int) ([]*Photo, error) {
	photos := make([]*Photo, 0, 10)

	//language=MySQL
	rows, err := st.db.Query("SELECT `id`, `user_id`, `path` FROM `photos` WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select photos from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		item := &Photo{}
		err := rows.Scan(&item.ID, &item.UserID, &item.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse result into Photo struct: %w", err)
		}
		photos = append(photos, item)
	}

	return photos, nil
}


