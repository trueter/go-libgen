package main

type Task struct {
    UUID string
    SourceURL string `json:"sourceURL"`
    Format string    `json:"format"`
    Email string     `json:"email"`
    Status string    `json:"status"`
}
