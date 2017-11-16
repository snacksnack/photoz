package views

import "html/template"

type View struct {
	Template *template.Template //pointer to a template
}

func NewView(files ...string) *View {
	//create new view object and parse all template files
	files = append(files, "/layouts/footer.gohtml")

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
	}
}
