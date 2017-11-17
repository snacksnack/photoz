package controllers

import (
	"fmt"
	"net/http"

	"../views"
)

type Users struct {
	NewView *views.View
}

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
	fmt.Fprintln(w, r.PostFormValue("email"))
	fmt.Fprintln(w, r.PostFormValue("password"))
	fmt.Fprintln(w, "Temporary placeholder - submitted form info received.")
}
