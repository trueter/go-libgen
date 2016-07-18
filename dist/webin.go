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
	fmt.Println(book_url)
	return err
}

func post_handler(w http.ResponseWriter, r *http.Request) {


    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    pong, err := client.Ping().Result()
    fmt.Println(pong, err)

    r.ParseForm()
    book_url := r.FormValue("url")
    err = addFetchTask( book_url )

    if err != nil {
        panic( err )
    }

}


func main() {
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/book", post_handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
