package main

import (
	"net/http"

	"../views"
	"github.com/gorilla/mux"
)

var (
	homeView    *views.View
	contactView *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := homeView.Template.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := contactView.Template.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	homeView = views.NewView("../views/home.gohtml")
	contactView = views.NewView("../views/contact.gohtml")

	/*var err error
	homeTemplate, err = template.ParseFiles(
		"../views/home.gohtml",
		"../views/layouts/footer.gohtml",
	)
	if err != nil {
		panic(err)
	}

	contactTemplate, err = template.ParseFiles(
		"../views/contact.gohtml",
		"../views/layouts/footer.gohtml",
	)
	if err != nil {
		panic(err)
	}*/

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", r)
}
