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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
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

func viewHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/view/"):]
	p, _ := loadOrder(name)
	fmt.Fprintf(w, "<h1>%s</h1><h4>%s</h4>", p.Name, p.Details)
}