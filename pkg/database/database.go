package database

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"goblog/pkg/logger"
	"time"
)

var DB *sql.DB

func Initialize()  {
	initDB()
	createTables()
}

func initDB() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "Landy552",
		Addr:                 "sh-cdb-iti0tmqw.sql.tencentcdb.com:60874",
		Net:                  "tcp",
		DBName:               "ciroy_maker",
		AllowNativePasswords: true,
	}

	fmt.Println(config.FormatDSN())

	DB, err = sql.Open("mysql", config.FormatDSN())

	logger.LogError(err)

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接，失败报错
	err = DB.Ping()
	logger.LogError(err)
}

func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
	id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
	title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
	body longtext COLLATE utf8mb4_unicode_ci
);`

	_, err := DB.Exec(createArticlesSQL)
	logger.LogError(err)
}