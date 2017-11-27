package controllers

import (
	"fmt"
	"net/http"

	"../models"
	"../views"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

type SignupForm struct {
	Email    string `schema:"email"` //struct tag for the json package
	Password string `schema:"password"`
	Name     string `schema:"name"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//creates Users controller
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

//GET signup form
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

//POST signup form
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		//not the best way to handle an error:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)

	//fmt.Fprintln(w, r.PostFormValue("email"))
	//fmt.Fprintln(w, r.PostFormValue("password"))
}

// POST login: verifies password/email combo and logs in user
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprint(w, "Invalid email address.")
		case models.ErrInvalidPassword:
			fmt.Fprint(w, "Invalid password provided")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	cookie := http.Cookie{
		Name:  "email",
		Value: user.Email,
	}
	http.SetCookie(w, &cookie)
	fmt.Fprint(w, user)
}
