package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir   string = "../views/layouts/"
	TemplateDir string = "../views/"
	TemplateExt string = ".gohtml"
)

type View struct {
	Template *template.Template //pointer to a template
	Layout   string
}

//Render is used to render current template with passed in data
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	// convert incoming data to type Data if necessary
	switch data.(type) {
	case Data:
		//do nothing
	default:
		data = Data{
			Yield: data,
		}
	}
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

//create new view object and parse all template files
func NewView(layout string, files ...string) *View {

	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

//takes slice of template filepath strings. prepends directory to file
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

//takes slice of template filepath strings. appends extension to file
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
