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
	db.AutoMigrate(&User{})

	//var u User
	//db.First(&u)
	// SELECT * FROM users ORDER BY id LIMIT 1;
	//fmt.Println(u)

	//db.First(&u, 4)
	//db.First(&u, "color = ?", "red")
	// SELECT * FROM users WHERE id = 4;
	//fmt.Println(u)

	//db.Last(&u)
	// SELECT * FROM users ORDER BY id DESC LIMIT 1;
	//fmt.Println(u)

	//db.Where("name = ? AND color = ?", "test2", "orange").First(&u)
	//SELECT * FROM users WHERE name = 'test2' AND color = 'orange' limit 1;
	//fmt.Println(u)

	//db.Where("name in (?)", []string{"test1", "test2"}).Find(&u)
	//SELECT * FROM "users" WHERE ((name in ('test1','test2')))
	//fmt.Println(u)

	//db.Where("color = ?", "blue").
	//	Where("id > ?", 3).
	//	First(&u)
	// SELECT * FROM "users" WHERE ((color = 'blue') AND (id > '3')) ORDER BY "users"."id" ASC LIMIT 1
	//fmt.Println(u)

	//var u User = User{
	//	Color: "red",
	//	Email: "test5@test.com",
	//}
	//query by resource
	//db.Where(u).First(&u)
	// SELECT * FROM "users" WHERE (("users"."email" = 'test5@test.com') AND ("users"."color" = 'red')) ORDER BY "users"."id" ASC LIMIT 1
	//fmt.Println(u)

	var users []User
	db.Find(&users)
	// SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL
	fmt.Println(len(users))
	fmt.Println(users)
}
