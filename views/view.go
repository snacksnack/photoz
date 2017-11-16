package views

import "html/template"

type View struct {
	Template *template.Template //pointer to a template
	Layout   string
}

func NewView(layout string, files ...string) *View {
	//create new view object and parse all template files
	files = append(files,
		"../views/layouts/bootstrap.gohtml",
		"../views/layouts/footer.gohtml",
	)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}
