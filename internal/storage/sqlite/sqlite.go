package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"urlShortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(StoragePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", StoragePath)
	if err != nil {
		return nil, fmt.Errorf("error in %s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
	    id INTEGER PRIMARY KEY,
	    alias TEXT NOT NULL UNIQUE,
	    url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("error in %s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("error in %s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave, alias string) error {
	const op = "storage.sqlite.SaveUrl"
	stmt, err := s.db.Prepare(`INSERT INTO url(url, alias) VALUES(?,?)`)
	if err != nil {
		return fmt.Errorf("error in %s: %w", op, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqlErr, ok := err.(sqlite3.Error); ok && sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("error in %s: %w", op, storage.ErrURLExist)
		}
		return fmt.Errorf("error in %s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"
	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias = ?`)
	if err != nil {
		return "", fmt.Errorf("error in %s: %w", op, err)
	}
	var res string
	err = stmt.QueryRow(alias).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("error in %s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("error in %s: %w", op, err)
	}
	return res, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.sqlite.DeleteUrl"
	stmt, err := s.db.Prepare(`DELETE FROM url WHERE alias = ?`)
	if err != nil {
		return fmt.Errorf("error in %s: %w", op, err)
	}
	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("error in %s: %w", op, err)
	}

	rowsDel, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error in %s: %w", op, err)
	}
	if rowsDel == 0 {
		return fmt.Errorf("error in %s: %w", op, storage.ErrURLNotFound)
	}

	return nil
}
