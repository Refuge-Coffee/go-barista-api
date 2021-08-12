package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Order struct {
	Name string
	Number string
	Details []byte
}

func (o *Order) save() error {
	filename := "data/" + o.Name + ".txt"
	return ioutil.WriteFile(filename, o.Details, 0600)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadOrder(name string) (*Order, error) {
	filename := "data/" + name + ".txt"
	details, err := ioutil.ReadFile(filename)
	if err != nil {
			return nil, err
	}
	return &Order{Name: name, Details: details}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, o *Order) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
    fn(w, r, m[2])
	}
}

func editHandler(w http.ResponseWriter, r *http.Request, name string) {
	p, err := loadOrder(name)
	if err != nil {
			p = &Order{Name: name}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, name string) {
	details := r.FormValue("details")
	p := &Order{Name: name, Details: []byte(details)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+name, http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request, name string) {
	p, err := loadOrder(name)
	if err != nil {
		http.Redirect(w, r, "/edit/"+name, http.StatusFound)
		return
}
  renderTemplate(w, "view", p)
}