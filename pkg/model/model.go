package model

import (
	"fmt"
	"goblog/pkg/config"
	"goblog/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error

	var (
		host     = config.GetString("database.mysql.host")
		port     = config.GetString("database.mysql.port")
		database = config.GetString("database.mysql.database")
		username = config.GetString("database.mysql.username")
		password = config.GetString("database.mysql.password")
		charset  = config.GetString("database.mysql.charset")
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		username, password, host, port, database, charset, true, "Local")

	config := mysql.New(mysql.Config{
		// DSN: "root:own3306@tcp(server.pca7.com:3306)/goblog?charset=utf8&parseTime=True&loc=Local",
		DSN: dsn,
	})

	DB, err = gorm.Open(config, &gorm.Config{})

	logger.LogError(err)
	return DB
}
