package url

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

func (s *Store) GetUrlsByUserID(id int64) ([]*types.UrlResponse, error) {
	rows, err := s.db.Query("SELECT * FROM urls WHERE userID = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []*types.UrlResponse
	for rows.Next() {
		url, err := scanRowIntoUrl(rows)
		if err != nil {
			return nil, err
		}
		urls = append(urls, &types.UrlResponse{ShortUrl: url.ShortUrl, FullUrl: url.FullUrl})
	}
	log.Println(urls, id)

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func (s *Store) InsertClickData(urlID int64, ip, country, device, platform, browser, language string) error {
	_, err := s.db.Exec("INSERT INTO clicks (urlID, clickedAt, ip_address, country, device, platform, browser, language) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", urlID, time.Now(), ip, country, device, platform, browser, language)
	if err != nil {
		return fmt.Errorf("error inserting click data: %w", err)
	}

	return nil
}

func (s *Store) GetClicksByID(id int64) ([]*types.Click, error) {
	rows, err := s.db.Query("SELECT * FROM clicks WHERE urlID = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clicks []*types.Click
	for rows.Next() {
		click, err := scanRowIntoClick(rows)
		if err != nil {
			return nil, err
		}
		clicks = append(clicks, click)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return clicks, nil
}

func (s *Store) GetClicksByHour(urlID int64) ([]types.HourlyClickStat, error) {
	rows, err := s.db.Query(`
        SELECT HOUR(clickedAt) AS hour, COUNT(*) AS click_count
        FROM clicks
        WHERE clickedAt >= NOW() - INTERVAL 1 DAY AND urlID = ?
        GROUP BY HOUR(clickedAt)
        ORDER BY HOUR(clickedAt)`, urlID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []types.HourlyClickStat
	for rows.Next() {
		var stat types.HourlyClickStat
		if err := rows.Scan(&stat.Hour, &stat.ClickCount); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}
	return stats, nil
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

func scanRowIntoClick(row *sql.Rows) (*types.Click, error) {
	click := new(types.Click)
	err := row.Scan(
		&click.ID,
		&click.UrlID,
		&click.Timestamp,
		&click.IPAddress,
		&click.Country,
		&click.Device,
		&click.Platform,
		&click.Browser,
		&click.Language,
	)
	if err != nil {
		return nil, err
	}
	return click, nil
}
