package main

import (
    "encoding/json"
    "github.com/satori/go.uuid"
    "html/template"
    "net/http"
    "log"
)

func serveHTTP( ) {
    // External API
    http.HandleFunc("/", indexHandler)
    // Book route handler
    http.HandleFunc("/save", saveHandler)

    // Serve the static page
    static := http.FileServer(http.Dir("static"))
    http.Handle( "/static/", http.StripPrefix("/static/", static))

    log.Fatal( http.ListenAndServe( ":3001", nil ))
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    // POST handler
    r.ParseForm()

    // Build task object
    task := Task {
        UUID     : uuid.NewV4().String(),
        SourceURL: r.FormValue("url"),
        Format   : r.FormValue("format"),
        Email    : r.FormValue("email"),
        Status   : "created"}

    // Add task to queue
    pushTask(task)
    //if err != nil {
    //    http.Error( w, err, http.StatusInternalServerError )
    //    return
    //}
    log.Print("Task added to queue:  %s\n", task)

    w.WriteHeader(200)
    w.Header().Set("Content-Type", "application/json")
    response, _ := json.Marshal(task)
    w.Write(response)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("form.html")
    t.Execute(w, nil)
}
