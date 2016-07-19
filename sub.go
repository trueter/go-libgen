package main

import (
	"fmt"
	//"log"
	"gopkg.in/redis.v4"
)

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Publish("tasks", "hello").Err()
	if err != nil {
		panic(err)
	}
	for {
		pubsub, err := client.Subscribe("tasks")
		if err != nil {
			panic(err)
		}

		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			panic(err)
		}

		fmt.Println(msg.Channel, msg.Payload)
	}

}
