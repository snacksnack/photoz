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
	fmt.Fprintln(w, "Temporary placeholder - submitted form info received.")
}
