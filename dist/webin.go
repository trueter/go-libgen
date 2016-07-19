package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "gopkg.in/redis.v4"
)


func indexHandler(w http.ResponseWriter, r *http.Request) {

    t, _ := template.ParseFiles("static/form.html")
    t.Execute(w, nil)
}


func addFetchTask(book_url string)(err error) {

    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    table:="toGet"

    resp, err := client.RPush(table,book_url).Result()
    fmt.Println("r:",resp)
    return err
}

func post_handler(w http.ResponseWriter, r *http.Request) {

    r.ParseForm()
    book_url := r.FormValue("url")
    err := addFetchTask( book_url )

    if err != nil {
        panic( err )
    }
   
}


func main() {

    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/book", post_handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
