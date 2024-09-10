package user

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

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email=?", email)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user with this mail doesn't exist")
	}
	return user, nil
}

func scanRowIntoUser(row *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByID(id int64) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user with this mail doesn't exist")
	}
	return user, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Query("INSERT INTO users (name, email, password) VALUES (?, ?, ?);", user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
