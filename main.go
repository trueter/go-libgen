package main

import (
    "encoding/json"
    "gopkg.in/redis.v4"
    "log"
)


// Global queue client
var queueClient * redis.Client

// Observer routine that subscribes to Redis "task" queue to poll for tasks
func observe() {
    // Connect to redis
    pubsub, err := queueClient.Subscribe("tasks")
    scream(err)
    // Go do your thing forever
    for {
        msg, err := pubsub.ReceiveMessage()
        scream(err)
        var task Task
        err = json.Unmarshal([]byte( msg.Payload ), &task)
        scream(err)
        switch task.Status {
            case "created":
              _, err = download(task)
            case "downloaded":
              _, err = convert(task)
            case "converted":
              _, err = mail(task)
            case "done":
              log.Print("Task Complete: ", task)
        }
        scream(err)
    }
}

// Redis task helpers
func pushTask(task Task) (err error) {
    log.Print("task", task)
    // Serialize the task object
    taskJSON, err := json.Marshal(task)
    if err != nil {
        return err
    }
    serializedTask := string(taskJSON)
    // Publish to redis task queue
    redisResult := queueClient.Publish("tasks", serializedTask)
    err = redisResult.Err()
    if err != nil {
        return err
    }
    log.Print("Result %s\n", redisResult.String())
    return nil
}

func main() {
    // Setup queue client
    queueClient = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    // Go subscribe to queue and poll for task
    go observe()
    // Server
    serveHTTP()
}
