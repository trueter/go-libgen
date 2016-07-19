package main

import (
    "bytes"
    "encoding/json"
    "io"
    "io/ioutil"
    "path/filepath"
    "fmt"
    "github.com/satori/go.uuid"
    "html/template"
    "log"
    "net/http"
    "net/url"
    "gopkg.in/redis.v4"
    "gopkg.in/gomail.v2"
    "os"
    "os/exec"
    "strings"
)


type Task struct {
    UUID string
    SourceURL string `json:"sourceURL"`
    Format string    `json:"format"`
    Email string     `json:"email"`
    Status string    `json:"status"`
}

var client * redis.Client


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

func pushTask( taskJSONPayload []byte ) ( err error ) {
    // fmt.Println("JSON ", taskJSONPayload )

    str := string( taskJSONPayload )

    redisResult := client.Publish( "tasks", str )

    err = redisResult.Err()
    if err != nil {
        return err
    }


    fmt.Println("Result %s\n", redisResult.String() )

    return nil
}

func post_handler(w http.ResponseWriter, r *http.Request) {

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

func mail( task Task, file * os.File ) {

    SMTP_HOST := "smtp.gmail.com"
    SENDER := "sender@example.org"
    USER := "bkcnvrt@gmail.com"
    PASSWORD := "PathAndTmpFile"
    body := "Download Link:\n"

    body = body + "http://torsten"

    /////
    m := gomail.NewMessage()
    m.SetHeader( "From", SENDER )
    m.SetHeader( "To", task.Email )
    m.SetHeader( "Subject", "Ebkcnvrt" )
    m.SetBody( "text/plain", body )
    m.Attach( file.Name() )

    d := gomail.NewPlainDialer( SMTP_HOST, 587, USER, PASSWORD )

    // Send the email to Bob, Cora and Dan.
    if err := d.DialAndSend( m ); err != nil {
        panic(err)
    }
}

func callback( task Task ) {

    // File: os.File

    file, err := ioutil.TempFile( "", "download" )
    if err != nil {
        panic( err )
    }



    // fmt.Println("Received ", taskPayloadJSON )
    fmt.Println("Temp File ", file.Name() )
    fmt.Println("Task", task )
    mail( task, file )

}

func observe() {

    pubsub, err := client.Subscribe("tasks")
    if err != nil {
        panic(err)
    }
    for {

        msg, err := pubsub.ReceiveMessage()

        if err != nil {
            panic(err)
        }

        var task Task
        err = json.Unmarshal( []byte( msg.Payload ), &task )
        if err != nil {
            panic(err)
        }

        callback( task )
    }
}

func main() {

    client = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    go observe()



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

    http.HandleFunc( "/", indexHandler )

    static := http.FileServer( http.Dir( "static" ) )

    http.Handle( "/static/", http.StripPrefix( "/static/", static ) )
    http.HandleFunc( "/save", post_handler )

    log.Fatal( http.ListenAndServe( ":3001", nil ) )
}


// http://unec.edu.az/application/uploads/2014/12/pdf-sample.pdf