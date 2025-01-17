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

func (s *Store) GetClicksByHourLast24Hours(urlID int64) ([]types.HourlyClickStat, error) {
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

func (s *Store) GetClicksByDayLast30Days(urlID int64) ([]types.DailyClickStat, error) {
	query := `
        SELECT 
		    DAY(clickedAt) AS day,
    		COUNT(*) AS click_count
		FROM 
    		clicks
		WHERE 
			urlID = ? AND
    		clickedAt >= NOW() - INTERVAL 30 DAY
		GROUP BY 
    		DAY(clickedAt)
		ORDER BY 
    		DAY(clickedAt);
    `

	rows, err := s.db.Query(query, urlID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dailyClickStats []types.DailyClickStat
	for rows.Next() {
		var stat types.DailyClickStat
		if err := rows.Scan(&stat.Day, &stat.ClickCount); err != nil {
			return nil, err
		}
		dailyClickStats = append(dailyClickStats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dailyClickStats, nil
}

func (s *Store) GetClicksByMonthLast12Months(urlID int64) ([]types.MonthlyClickStat, error) {
	query := `
        SELECT 
    		DATE_FORMAT(clickedAt, '%m') AS month, 
    		COUNT(*) AS click_count
		FROM 
    		clicks
		WHERE 
		    urlID = ? AND
			clickedAt >= NOW() - INTERVAL 12 MONTH 
		GROUP BY 
    		DATE_FORMAT(clickedAt, '%m')
		ORDER BY 
    		DATE_FORMAT(clickedAt, '%m');
    `

	rows, err := s.db.Query(query, urlID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monthlyClickStats []types.MonthlyClickStat
	for rows.Next() {
		var stat types.MonthlyClickStat
		if err := rows.Scan(&stat.Month, &stat.ClickCount); err != nil {
			return nil, err
		}
		monthlyClickStats = append(monthlyClickStats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return monthlyClickStats, nil
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
