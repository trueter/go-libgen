package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func get_book(book_url string){
    fmt.Printf("trying to get book %s", book_url)
}

func post_handler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    book_url := r.FormValue("book_url")
    fmt.Fprintf(w, "path:%s<br>", r.URL.Path[1:])
    fmt.Fprintf(w, "book_url:%s<br>", book_url)
    // validate res
    get_book(book_url)
}

func main() {

    //book_url:="";
    http.HandleFunc("/", handler)
    http.HandleFunc("/book", post_handler)
    http.ListenAndServe(":8080", nil)
}
