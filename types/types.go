package types

import "time"

type RequestUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SaveUrlPayload struct {
	UserID   int64  `json:"userID"`
	FullUrl  string `json:"fullUrl"`
	ShortUrl string `json:"shortUrl"`
}

type ModifyUrlPayload struct {
	FullUrl  string `json:"fullUrl"`
	ShortUrl string `json:"shortUrl"`
}

type RemoveUrlPayload struct {
	ShortUrl string `json:"shortUrl"`
}

type GetAllUrlPayload struct {
	UserID int64 `json:"userID"`
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt:"`
}

type Url struct {
	ID       int    `json:"id"`
	ShortUrl string `json:"shortUrl"`
	FullUrl  string `json:"fullUrl"`
	UserID   int64  `json:"userID"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

type UrlStore interface {
	GetUrlByShortUrl(shortUrl string) (*Url, error)
	CreateUrl(Url) error
	ModifyUrl(oldUrl Url, newUrl Url) error
	RemoveUrl(shortUrl string) error
	GetUrlsByUserID(id int) ([]*Url, error)
}
