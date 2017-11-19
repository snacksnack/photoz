package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "photoz"
	password = "photoz"
	dbname   = "photoz"
)

type User struct {
	gorm.Model //not inheritance - embedding
	Name       string
	Email      string `gorm:"not null; unique_index"`
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.DB().Ping(); err != nil {
		panic(err)
	}

	//db.DropTableIfExists(&User{}) -- NEVER DO THIS IN PRODUCTION, LOSER
	db.AutoMigrate(&User{})

	/*
		user := User{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Now(),
			},
		}
		fmt.Println(user.CreatedAt) // same as user.Model.CreatedAt
	*/
}
