package main

import (
    "io"
    "io/ioutil"
    "path/filepath"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "net/url"
    "os"
    "strings"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

    t, _ := template.ParseFiles("form.html")
    t.Execute(w, nil)
}

func downloadInto( file * os.File , url string) (err error) {

  // Get the data
  resp, err := http.Get(url)
  if err != nil {
    fmt.Printf( "err", err )
    return err
  }
  // Uncomment to display body
  // fmt.Printf( "resp", resp , "\n")
  defer resp.Body.Close()

  // Writer the body to file
  _, err = io.Copy(file, resp.Body)
  if err != nil  {
    return err
  }

  return nil
}

func getFileNameFromURL( _url string ) string {

    u, err := url.Parse( _url )
    if err != nil {
        panic( err )
    }

    str := u.Path

    i := strings.LastIndex( str, "/" )
    return str[ i + 1 : len( str ) ]

}

func getBook( book_url string ) ( err error, fileName string, filePath string ) {

    fileName = getFileNameFromURL( book_url )

    fmt.Printf("filenam is %s\n", fileName )

    file, err := ioutil.TempFile( "", "download" )
    if err != nil {
        return err, "", ""
    }


    err2 := downloadInto( file, book_url )
    if err2 != nil {
        return err2, "", ""
    }

    filePath, err3 := filepath.Abs( file.Name() )
    if err3 != nil {
        return err3, "", ""
    }

    fmt.Println("Prepared ", filePath )
    return nil, fileName, filePath


    // os.Remove(file.Name())
}


func post_handler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    book_url := r.FormValue("url")
//    fmt.Fprintf(w, "path:%s<br>", r.URL.Path[1:])
    fmt.Fprintf(w, "%s", book_url)
    // validate res
    getBook(book_url)
}



func main() {
    http.HandleFunc("/", indexHandler)
    static := http.FileServer( http.Dir( "static" ) )
    http.Handle("/static/", http.StripPrefix("/static/", static))
    http.HandleFunc("/book", post_handler)
    log.Fatal(http.ListenAndServe(":8080", nil))

    err, fileName, filePath := getBook( "http://www.orimi.com/pdf-test.pdf" )
    if err != nil {
        panic( err )
    }

    fmt.Println("fileName ", fileName )
    fmt.Println("filePath ", filePath )

}
