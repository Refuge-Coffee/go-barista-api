package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
)

type Order struct {
	Name string
	Number string
	Details []byte
}

func (o *Order) save() error {
	filename := o.Name + ".txt"
	return ioutil.WriteFile(filename, o.Details, 0600)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadOrder(name string) (*Order, error) {
	filename := name + ".txt"
	details, err := ioutil.ReadFile(filename)
	if err != nil {
			return nil, err
	}
	return &Order{Name: name, Details: details}, nil
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/edit/"):]
	p, err := loadOrder(name)
	if err != nil {
			p = &Order{Name: name}
	}
	fmt.Fprintf(w, "<h1>Editing %s</h1>"+
			"<form action=\"/save/%s\" method=\"POST\">"+
			"<textarea name=\"details\">%s</textarea><br>"+
			"<input type=\"submit\" value=\"Save\">"+
			"</form>",
			p.Name, p.Name, p.Details)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/save/"):]
	details := r.FormValue("details")
	p := &Order{Name: name, Details: []byte(details)}
	p.save()
	http.Redirect(w, r, "/view/"+name, http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/view/"):]
	p, _ := loadOrder(name)
	fmt.Fprintf(w, "<h1>%s</h1><h4>%s</h4>", p.Name, p.Details)
}