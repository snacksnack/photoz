package main

import "fmt"

type Dog struct{}

type Cat struct{}

func (c Cat) Speak() {
	fmt.Println("meow")
}

func (d Dog) Speak() {
	fmt.Println("woof")
}

type Husky struct {
	//dog Dog
	//Dog //now you Husky can call any method available to Dog type
	Speaker
}

type Speaker interface {
	Speak()
}

type SpeakerPrefixer struct {
	Speaker
}

func (sp SpeakerPrefixer) Speak() {
	fmt.Print("Prefix: ")
	sp.Speaker.Speak()
}

/*
type UserReader interface {
	ByID(id uint) (*User, error)
}
*/

func main() {
	//h := Husky{Dog{}}
	//h := Husky{Cat{}}
	h := Husky{SpeakerPrefixer{Cat{}}}
	h.Speak() //equal to h.Dog.Speak()
}
