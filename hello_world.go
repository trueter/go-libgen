package main

import (
    "bytes"
    "encoding/json"
    "io"
    "github.com/satori/go.uuid"
    "html/template"
    "log"
    "net/http"
    "net/url"
    "gopkg.in/redis.v4"
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
    FilePath string  `json:"filePath"`
}

var client * redis.Client

func indexHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("form.html")
    t.Execute(w, nil)
}


func mail(task Task) (err error) {
  return nil
}


func download(task Task) (err error) {
    log.Print("Download Task Recieved: ", task)
   // Create temp file
    tmpDir := os.TempDir()
    filePath := tmpDir + task.UUID

    file, err := os.Create(filePath)
    scream(err)
    log.Print("Created Temp File at ", file.Name())

    // Get the data
    resp, err := http.Get(task.SourceURL)
    scream(err)
    defer resp.Body.Close()

    // Writer the body to file
    _, err = io.Copy(file, resp.Body)
    scream(err)

    // Update the task
    task.Status = "downloaded"
    log.Print("Download Task Complete: ", task)
    // Update task queue
    pushTask(task)
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

func convert(task Task) (err error) {
    // Assuming temp files are in tmp
    tmpDir := os.TempDir()
    tmpFile := tmpDir + task.UUID
    pdfFile := tmpDir + task.UUID + ".pdf"
    convertedFile := tmpDir + task.UUID + "." + task.Format
    log.Print("tmp File:", tmpFile)

    // Add PDF extension
    err = os.Rename(tmpFile, pdfFile)
    scream(err)
    log.Print("PDF File: %s", pdfFile)


    cmd := exec.Command("ebook-convert", pdfFile, convertedFile)

    var out bytes.Buffer
    cmd.Stdout = &out
    err = cmd.Run()
    scream(err)
    log.Print("%s File: %s", convertedFile)
    //    log.Print("out:%s",out.String()) ??

    // Update the task
    task.Status = "converted"
    pushTask(task)

    return nil
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


func observe() {
    pubsub, err := client.Subscribe("tasks")
    scream(err)
    // Go do your thing forever
    for {
        msg, err := pubsub.ReceiveMessage()
        scream(err)
        var task Task
        err = json.Unmarshal( []byte( msg.Payload ), &task )
        scream(err)
        switch task.Status {
          case "created":
            download(task)
          case "downloaded":
            convert(task)
          case "converted":
            mail(task)
        }
    }
}



// Util functions
func scream(err error) {
    if err != nil {
      log.Fatal(err)
      panic(err)
    }
}

// Redis task helpers
func pushTask(task Task) (err error) {
    taskJSON, err := json.Marshal(task)
    scream(err)
    serializedTask := string(taskJSON)

    redisResult := client.Publish("tasks", serializedTask)
    err = redisResult.Err()
    scream(err)

    log.Print("Result %s\n", redisResult.String())
    return nil
}


func main() {

    client = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    // Set subscribe in a goroutine to poll for tasks
    go observe()

    // External API
    http.HandleFunc("/", indexHandler)
    // Book route handler
    http.HandleFunc("/save", saveHandler)

    // Serve the static page
    static := http.FileServer(http.Dir("static"))
    http.Handle( "/static/", http.StripPrefix( "/static/", static ) )

    log.Fatal( http.ListenAndServe( ":3001", nil ) )
}


// http://unec.edu.az/application/uploads/2014/12/pdf-sample.pdf
