package url

import (
	"database/sql"

	"github.com/nirav114/url-shortner-backend.git/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) GetUrlByShortUrl(shortUrl string) (*types.Url, error) {
	return nil, nil
}

func (s *Store) CreateUrl(url types.Url) error {
	return nil
}

func (s *Store) ModifyUrl(oldUrl types.Url, newUrl types.Url) error {
	return nil
}

func (s *Store) RemoveUrl(shortUrl string) error {
	return nil
}

func (s *Store) GetUrlsByUserID(id int) ([]*types.Url, error) {
	return []*types.Url{}, nil
}
