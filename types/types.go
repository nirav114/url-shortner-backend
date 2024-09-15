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

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt:"`
}

type Url struct {
	ID       int64  `json:"id"`
	ShortUrl string `json:"shortUrl"`
	FullUrl  string `json:"fullUrl"`
	UserID   int64  `json:"userID"`
}

type UrlResponse struct {
	ShortUrl string `json:"shortUrl"`
	FullUrl  string `json:"fullUrl"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int64) (*User, error)
	CreateUser(User) error
}

type UrlStore interface {
	GetUrlByShortUrl(shortUrl string) (*Url, error)
	CreateUrl(Url) error
	ModifyUrl(oldUrl Url, newUrl Url) error
	RemoveUrl(shortUrl string) error
	GetUrlsByUserID(id int64) ([]*UrlResponse, error)
	InsertClickData(urlID int64, ip, country, device, platform, browser, language string) error
}

type UserOTPData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	OTP      string `json:"otp"`
}

type VerifyOTPPayload struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}
