package main

import (
    "encoding/json"
    "io"
    "io/ioutil"
    "path/filepath"
    "bytes"
    "fmt"
    "github.com/satori/go.uuid"
    "html/template"
    "log"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "strings"
)


type Task struct {
    UUID string
    sourceURL string
    format string
    email string
    status string
}


func indexHandler(w http.ResponseWriter, r *http.Request) {

    t, _ := template.ParseFiles("form.html")
    t.Execute(w, nil)
}



func download( file * os.File , url string) (err error) {

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

func convertBook(tmpPathAndFileName string, outFilename string, format string) (err error, outPathAndFileName string) {

    i := strings.LastIndex( tmpPathAndFileName, "/" )
    tmpPath := tmpPathAndFileName[0:i]
    log.Print("tmp Path:%s",tmpPath)

    err = os.Rename(tmpPathAndFileName, tmpPathAndFileName + ".pdf")
    tmpPathAndFileName = tmpPathAndFileName + ".pdf"

    outPathAndFilename := tmpPath + "/" + outFilename + format
    log.Print("tmpPathAndFileName:%s",tmpPathAndFileName)
    log.Print("outPathAndFilename:%s",outPathAndFilename)


    if err != nil{
        return err, ""
    }

    cmd := exec.Command("ebook-convert", tmpPathAndFileName, outPathAndFilename)
    var out bytes.Buffer
    cmd.Stdout = &out
    err = cmd.Run()
    if err != nil{
        return err, ""
    }
    log.Print("out:%s",out.String())

    return nil, outPathAndFileName
}


// func cleanAndDisarm( inPathAndFilename string ) string {
//     charsToRemove := "'\"\\`{[(<>)]}^|!?#$*%&=$ ;,:."
//     inPathAndFilename.translate( None, charsToRemove )
//     return inPathAndFilename
// }


func getBook( book_url string ) ( err error, fileName string, filePath string ) {

    fileName = getFileNameFromURL( book_url )

    fmt.Printf("filenam is %s\n", fileName )

    file, err := ioutil.TempFile( "", "download" )
    if err != nil {
        return err, "", ""
    }


    err2 := download( file, book_url )
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

func pushTask( task Task ) ( err error ) {
    return nil
}

func post_handler(w http.ResponseWriter, r *http.Request) {

    r.ParseForm()

    task := Task{
        UUID     : uuid.NewV4().String(),
        sourceURL: r.FormValue("url"),
        format   : r.FormValue("format"),
        email    : r.FormValue("email"),
        status   : "created"}

    err := pushTask( task )
    if err != nil {
        http.Error( w, err.Error(), http.StatusInternalServerError )
        return
    }

    responseJSON, err := json.Marshal( task )
    if err != nil {
        http.Error( w, err.Error(), http.StatusInternalServerError )
        return
    }


    w.WriteHeader( 200 )
    w.Header().Set( "Content-Type", "application/json" )
    w.Write( responseJSON )

    // format := ".mobi"

    // //
    // err, outFilename, filePath := getBook( "http://www.orimi.com/pdf-test.pdf" )
    // if err != nil {
    //     panic( err )
    // }

    // // filePath = cleanAndDisarm( filePath )

    // outFilename = outFilename
    // fmt.Println("outFilename ", outFilename )
    // fmt.Println("filePath ", filePath )
    // //

    // err, outPathAndFilename := convertBook(filePath, outFilename, format)
    // fmt.Println("outPathAndFilename ", outPathAndFilename )

}



func main() {



    // format := ".mobi"

    // err, outFilename, filePath := getBook( "http://www.orimi.com/pdf-test.pdf" )
    // if err != nil {
    //     panic( err )
    // }

    // // filePath = cleanAndDisarm( filePath )

    // outFilename = outFilename
    // fmt.Println("outFilename ", outFilename )
    // fmt.Println("filePath ", filePath )
    // //

    // err, outPathAndFilename := convertBook(filePath, outFilename, format)
    // fmt.Println("outPathAndFilename ", outPathAndFilename )
    // if err != nil {
    //     panic( err )
    // }

    http.HandleFunc("/", indexHandler)
    static := http.FileServer( http.Dir( "static" ) )
    http.Handle("/static/", http.StripPrefix("/static/", static))
    http.HandleFunc("/book", post_handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
