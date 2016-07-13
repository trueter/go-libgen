package main

import (
	//"fmt"
    "html/template"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("form.html")
	t.Execute(w, "hi")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
