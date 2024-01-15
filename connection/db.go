package connection

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

var DB_HOST string = GoDotEnvVariable("DB_HOST")
var DB_DATABASE string = GoDotEnvVariable("DB_DATABASE")
var DB_USERNAME string = GoDotEnvVariable("DB_USERNAME")
var DB_PASSWORD string = GoDotEnvVariable("DB_PASSWORD")
var DB_PORT string = GoDotEnvVariable("DB_PORT")

func Connect() {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USERNAME, DB_PASSWORD, DB_HOST, DB_PORT, DB_DATABASE)
	d, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db = d
}

func GetDb() *gorm.DB {
	return db
}
