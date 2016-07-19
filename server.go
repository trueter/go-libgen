package main

import (
    "encoding/json"
    "fmt"
    "github.com/satori/go.uuid"
    "html/template"
    "net/http"
    "log"
)


func serveHTTP( ) {

    http.HandleFunc( "/", IndexHandleFunc )

    static := http.FileServer( http.Dir( "static" ) )

    http.Handle( "/static/", http.StripPrefix( "/static/", static ) )
    http.HandleFunc( "/save", PostHandleFunc )

    log.Fatal( http.ListenAndServe( ":3001", nil ) )
}


func PostHandleFunc( w http.ResponseWriter, r *http.Request ) {

    r.ParseForm()

    task := &Task{
        UUID     : uuid.NewV4().String(),
        SourceURL: r.FormValue( "url" ),
        Format   : r.FormValue( "format" ),
        Email    : r.FormValue( "email" ),
        Status   : "created"}

    fmt.Println( "Task %s\n", task )

    responseJSON, err := json.Marshal( task )
    if err != nil {
        http.Error( w, err.Error(), http.StatusInternalServerError )
        return
    }


    err = pushTask( responseJSON )
    if err != nil {
        http.Error( w, err.Error(), http.StatusInternalServerError )
        return
    }


    w.WriteHeader( 200 )
    w.Header().Set( "Content-Type", "application/json" )
    w.Write( responseJSON )
}

func IndexHandleFunc(w http.ResponseWriter, r *http.Request) {

    t, _ := template.ParseFiles("form.html")
    t.Execute(w, nil)
}
