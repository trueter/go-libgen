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

    static := http.FileServer( http.Dir( "static" ) )
    http.Handle("/static/", http.StripPrefix("/static/", static))

	http.ListenAndServe(":8080", nil)
}
