package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	Color      string
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

	db.LogMode(true)

	//db.DropTableIfExists(&User{}) -- NEVER DO THIS IN PRODUCTION, LOSER
	db.AutoMigrate(&User{})

	name, email, color := getInfo()
	u := User{
		Name:  name,
		Email: email,
		Color: color,
	}
	if err = db.Create(&u).Error; err != nil {
		panic(err)
	} else {
		fmt.Printf("%v+\n", u)
	}
}

func getInfo() (name, email, color string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter your name:")
	name, _ = reader.ReadString('\n')
	fmt.Println("Please enter your email address:")
	email, _ = reader.ReadString('\n')
	fmt.Println("Please enter your favorite color:")
	color, _ = reader.ReadString('\n')
	color = strings.TrimSpace(color)
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	return name, email, color
}
