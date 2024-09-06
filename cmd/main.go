package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/nirav114/url-shortner-backend.git/cmd/api"
	"github.com/nirav114/url-shortner-backend.git/config"
	"github.com/nirav114/url-shortner-backend.git/db"
)

func main() {
	db, err := db.NewMySqlStorage(mysql.Config{
		User:                 config.EnvConfig.DBUser,
		Passwd:               config.EnvConfig.DBPassword,
		Addr:                 config.EnvConfig.DBAddress,
		DBName:               config.EnvConfig.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)

	server := api.NewApiServer(":3000", nil)
	if err := server.Run(); err != nil {
		log.Fatal()
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database successfully!")
}
