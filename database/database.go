package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Models   []interface{}
}

var dialector gorm.Dialector
var DB *gorm.DB

func Connect(config DBConfig) {

	switch config.Driver {
	case "sqlite":
		dialector = sqlite.Open(config.Database)
		break
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to open database connection: ", err)
	}

	Migrator.migrate()

	////loop through models and create tables
	//for _, model := range config.Models {
	//	err = DB.AutoMigrate(model)
	//	if err != nil {
	//		log.Fatalln("Failed to create table: ", err)
	//	}
	//}

}
