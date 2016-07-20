package main

import (
    "bytes"
    "gopkg.in/gomail.v2"
    "net/http"
    "io"
    "log"
    "os"
    "os/exec"
)

// Task methods
func mail(task Task) (result Task, err error) {
    // Setup the mailer
    SMTP_HOST := "smtp.gmail.com"
    SENDER := os.Getenv("SENDER")
    USER := os.Getenv("USER")
    PASSWORD := os.Getenv("PASSWORD")
    mailDialer := gomail.NewPlainDialer(SMTP_HOST, 587, USER, PASSWORD)
    // Get the file
    tmpDir := os.TempDir()
    filePath := tmpDir + task.UUID
    fileName := filePath + "." + task.Format
    file, err := os.Open(fileName)
    if err != nil {
        return task, err
    }
    // Build mail
    body := "Download Link:\n"
    message := gomail.NewMessage()
    message.SetHeader("From", SENDER)
    message.SetHeader("To", task.Email)
    message.SetHeader("Subject", "Ebkcnvrt")
    message.SetBody("text/plain", body)
    message.Attach(file.Name())
    // Sendmail
    err = mailDialer.DialAndSend(message)
    if err != nil {
        return task, err
    }
    // Update task
    task.Status = "done"
    log.Print("Mail task complete: ", task)
    err = pushTask(task)
    if err != nil {
      return task, err
    }
    return task, nil
}

func download(task Task) (result Task, err error) {
    log.Print("Download Task Recieved: ", task)
    // Create temp file
    tmpDir := os.TempDir()
    filePath := tmpDir + task.UUID
    file, err := os.Create(filePath)
    if err != nil {
        return task, err
    }
    log.Print("Created Temp File at ", file.Name())
    // Get the data
    resp, err := http.Get(task.SourceURL)
    if err != nil {
        return task, err
    }
    defer resp.Body.Close()
    // Write the body to file
    _, err = io.Copy(file, resp.Body)
    scream(err)
    // Update the task
    task.Status = "downloaded"
    log.Print("Download task complete: ", task)
    err = pushTask(task)
    if err != nil {
      return task, err
    }
    return task, nil
}

func convert(task Task) (result Task, err error) {
    // Assuming temp files are in tmp
    tmpDir := os.TempDir()
    tmpFile := tmpDir + task.UUID
    pdfFile := tmpDir + task.UUID + ".pdf"
    convertedFile := tmpDir + task.UUID + "." + task.Format
    log.Print("tmp File:", tmpFile)
    // Add PDF extension
    err = os.Rename(tmpFile, pdfFile)
    if err != nil {
        return task, err
    }
    log.Print("PDF File: %s", pdfFile)
    // Convert
    cmd := exec.Command("ebook-convert", pdfFile, convertedFile)
    var out bytes.Buffer
    cmd.Stdout = &out
    err = cmd.Run()
    if err != nil {
        return task, err
    }
    log.Print("%s File: %s", convertedFile)
    // Update the task
    task.Status = "converted"
    log.Print("Convert task complete: ", task)
    err = pushTask(task)
    if err != nil {
      return task, err
    }
    return task, nil
}
