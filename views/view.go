package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"

	"photoz/context"
)

var (
	LayoutDir   string = "./views/layouts/"
	TemplateDir string = "./views/"
	TemplateExt string = ".gohtml"
)

type View struct {
	Template *template.Template //pointer to a template
	Layout   string
}

//Render is used to render current template with passed in data
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	// convert incoming data to type Data if necessary
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	if alert := getAlert(r); alert != nil {
		vd.Alert = alert
		clearAlert(w)
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. Please try again.", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

//create new view object and parse all template files
func NewView(layout string, files ...string) *View {

	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)

	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("csrfField is not implemented")
		},
	}).ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
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
