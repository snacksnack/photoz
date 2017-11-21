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
	Color      string
	Orders     []Order //orders will not be pull by default - must tell gorm to preload
}

type Order struct {
	gorm.Model
	UserId      uint
	Amount      int
	Description string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(true)
	db.AutoMigrate(&User{}, &Order{})

	var u User
	if err := db.Preload("Orders").First(&u).Error; err != nil {
		panic(err)
	}
	fmt.Println(u)
	fmt.Println(len(u.Orders))
	//SELECT * FROM "orders"  WHERE "orders"."deleted_at" IS NULL AND (("user_id" IN('2'))) ORDER BY "orders"."id" ASC

	var users []User
	if err := db.Preload("Orders").Find(&users).Error; err != nil {
		panic(err)
	}
	fmt.Println(users)
	//SELECT * FROM "orders"  WHERE "orders"."deleted_at" IS NULL AND (("user_id" IN ('2','4','5','6','7')))

	//createOrder(db, u, 1001, "pants1")
	//createOrder(db, u, 1002, "pants2")
	//createOrder(db, u, 1003, "pants3")
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	err := db.Create(&Order{
		UserId:      user.ID,
		Amount:      amount,
		Description: desc,
	}).Error
	if err != nil {
		panic(err)
	}
}
