package views

import (
	"fmt"
	"html/template"
)

type View struct {
	Template *template.Template //pointer to a template
}

func NewView(files ...string) *View {
	//create new view object and parse all template files
	fmt.Print(files)
	files = append(files, "../views/layouts/footer.gohtml")
	fmt.Print(files)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
	}
}
