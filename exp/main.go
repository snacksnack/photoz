package main

import (
	"html/template"
	"os"
)

type User struct {
	Name          string
	Dog           Dog
	Age           int
	Temp          float64
	Colors        []string
	Personalities map[string]string
}

type Dog struct {
	Name string
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	data := User{
		Name: "Barry Larry",
		Dog: Dog{
			Name: "Larry Jr.",
		},
		Age:    3,
		Temp:   98.6,
		Colors: []string{"orange", "green", "blue", "red"},
		Personalities: map[string]string{
			"bob":  "crazy",
			"gary": "crazier",
		},
	}

	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}

	data.Name = "<script>alert('hello there.')</script>"
	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
