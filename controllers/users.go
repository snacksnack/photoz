package controllers

import (
	"fmt"
	"net/http"

	"../views"
	"github.com/gorilla/schema"
)

type Users struct {
	NewView *views.View
}

type SignupForm struct {
	Email    string `schema:"email"` //struct tag for the json package
	Password string `schema:"password"`
}

//creates Users controller
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "../views/users/new.gohtml"),
	}
}

//GET signup form
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

//POST signup form
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	dec := schema.NewDecoder()
	var form SignupForm
	if err := dec.Decode(&form, r.PostForm); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)

	//fmt.Fprintln(w, r.PostFormValue("email"))
	//fmt.Fprintln(w, r.PostFormValue("password"))
}
