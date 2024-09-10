package url

import (
	"database/sql"
	"fmt"

	"github.com/nirav114/url-shortner-backend.git/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) GetUrlByShortUrl(shortUrl string) (*types.Url, error) {
	rows, err := s.db.Query("SELECT * FROM urls WHERE shortUrl=?", shortUrl)
	if err != nil {
		return nil, err
	}

	url := new(types.Url)
	for rows.Next() {
		url, err = scanRowIntoUrl(rows)
		if err != nil {
			return nil, err
		}
	}

	if url.ID == 0 {
		return nil, fmt.Errorf("url with this shortUrl doesn't exist")
	}
	return url, nil
}

func (s *Store) CreateUrl(url types.Url) error {
	_, err := s.db.Exec("INSERT INTO urls (shortUrl, fullUrl, userID) VALUES (?, ?, ?);", url.ShortUrl, url.FullUrl, url.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ModifyUrl(oldUrl types.Url, newUrl types.Url) error {
	_, err := s.db.Exec("UPDATE urls SET fullUrl = ? WHERE shortUrl = ?", newUrl.FullUrl, newUrl.ShortUrl)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) RemoveUrl(shortUrl string) error {
	_, err := s.db.Exec("DELETE FROM urls WHERE shortUrl = ?", shortUrl)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetUrlsByUserID(id int) ([]*types.Url, error) {
	return []*types.Url{}, nil
}

func scanRowIntoUrl(row *sql.Rows) (*types.Url, error) {
	url := new(types.Url)
	err := row.Scan(
		&url.ID,
		&url.ShortUrl,
		&url.FullUrl,
		&url.UserID,
	)
	if err != nil {
		return nil, err
	}
	return url, nil
}
