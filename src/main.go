package main

import (
	"fmt"
	"net/http"

	"../controllers"
	"../models"
	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "photoz"
	password = "photoz"
	dbname   = "photoz_test"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	must(err)
	defer us.Close()
	//us.DestructiveReset()
	us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	//look into changing this to Handle that simply calls NewView template
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	//use Handle as we're not calling a function here - calling template instead
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
