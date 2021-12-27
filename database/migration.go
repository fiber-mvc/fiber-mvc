package database

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Migration struct {
	ID         uint   `gorm:"primarykey"`
	Name       string `gorm:"type:varchar(255);"`
	ExecutedAt time.Time
}
type _Migrator struct{}

var Migrator _Migrator = _Migrator{}

func getSqlFiles() []string {

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(workingDir + "/migrations")

	if err != nil {
		log.Fatal(err)
	}

	//loop through files
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, workingDir+"/migrations/"+file.Name())
	}

	return fileNames
}

func (m *_Migrator) migrate() {
	DB.AutoMigrate(&Migration{})
	files := getSqlFiles()

	for _, file := range files {

		//check if migration has already been executed
		var migration Migration
		DB.Where("name = ?", file).First(&migration)

		if migration.ID == 0 {

			log.Println("Migrating: " + file)

			//read the file content to a string
			content, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}

			DB.Exec(string(content))

			//save migration
			migration := Migration{Name: file}
			migration.ExecutedAt = time.Now()
			DB.Create(&migration)

		}
	}

}
