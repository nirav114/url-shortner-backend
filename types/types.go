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

type GetStatsPayload struct {
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
	GetClicksByID(urlID int64) ([]*Click, error)
	GetClicksByHourLast24Hours(urlID int64) ([]HourlyClickStat, error)
	GetClicksByDayLast30Days(urlID int64) ([]DailyClickStat, error)
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

type Click struct {
	ID        int64     `json:"id"`
	UrlID     int64     `json:"urlID"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip_address"`
	Country   string    `json:"country"`
	Device    string    `json:"device"`
	Platform  string    `json:"platform"`
	Browser   string    `json:"browser"`
	Language  string    `json:"language"`
}

type HourlyClickStat struct {
	Hour       int `json:"hour"`
	ClickCount int `json:"click_count"`
}

type DailyClickStat struct {
	Day        string `json:"day"`
	ClickCount int    `json:"click_count"`
}
