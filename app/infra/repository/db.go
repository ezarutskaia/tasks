package repository

import (
	"log"
	"gorm.io/gorm"
    "gorm.io/driver/mysql"
)

func SqlConnection() **gorm.DB {
	dsn := "root:secret@tcp(0.0.0.0:4306)/tasks?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
        log.Fatal(err)
    }
	return &db
}