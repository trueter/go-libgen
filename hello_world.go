package main

import (
    "fmt"
    "html/template"
	"log"
    "net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("form.html")
	t.Execute(w, "hi")
}


func getBook(book_url string){
    fmt.Printf("trying to get book %s", book_url)
}

func post_handler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    book_url := r.FormValue("book_url")
    fmt.Fprintf(w, "path:%s<br>", r.URL.Path[1:])
    fmt.Fprintf(w, "book_url:%s<br>", book_url)
    // validate res
    getBook(book_url)
}



func main() {
    http.HandleFunc("/", indexHandler)
    static := http.FileServer( http.Dir( "static" ) )
    http.Handle("/static/", http.StripPrefix("/static/", static))
    http.HandleFunc("/book", post_handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


