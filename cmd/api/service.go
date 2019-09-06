package main

import (
	"github.com/streadway/amqp"
	"bucket_file/cmd/api/controller"
	"log"
	_ "fmt"
	_ "bucket_file/constant"
	"bucket_file/model"
	_ "github.com/gin-gonic/gin"
	_ "net/http"
	_ "path/filepath"

	// "encoding/json"
	// "reflect"
	// "bucket_file/constant"

	// "bytes"
    // "net/http"
 
)

func Service() {
	err := controller.NewRouter(container)
	if err != nil {
		log.Fatal(err)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Fatalf("%s: %s", msg, err)
	}
}

func failOnErrorjson(err model.JsonErrorResponse, msg string) {
	log.Fatalf("%s: %s", msg, err)
}

func B2S(bs []uint8) string {
    ba := []byte{}
    for _, b := range bs[8:40] {
        ba = append(ba, byte(b))
    }
    return string(ba)
}

func CreateFile() {
	conn, err := amqp.Dial("amqp://test:test@192.168.30.94:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
	"upload file", // name
	false,   	 // durable
	false,   	 // delete when usused
	false,   	 // exclusive
	false,   	 // no-wait
	nil,    	 // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )
	  failOnError(err, "Failed to register a consumer")
	
	stopChan := make(chan bool)
	go func() {
		for d := range msgs {
			body := B2S(d.Body)
			_, code := container.MongoClient.File().FileRequest("PUT", "hunghunghung12", body + ".jpg", "../cloud_storage2/saved/" + body + ".jpg")
			if code != 200 {
				return 
			}
		}	
	}()
	<-stopChan
}

func createAdmin() {
	container.MongoClient.User().CreateAdmin()
}



// func RabitmqRecieve(routing_key string) (<-chan amqp.Delivery, error) {
// 	conn, err := amqp.Dial("amqp://test:test@192.168.30.94:5672/")
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	fmt.Println("2")
// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	fmt.Println("3")
// 	q, err := ch.QueueDeclare(
// 	routing_key, // name
// 	false,   	 // durable
// 	false,   	 // delete when usused
// 	false,   	 // exclusive
// 	false,   	 // no-wait
// 	nil,    	 // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")

// 	fmt.Println("4")
// 	msgs, err := ch.Consume(
// 		q.Name, // queue
// 		"",     // consumer
// 		true,   // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // args
// 	  )
// 	  fmt.Println(5)
// 	  failOnError(err, "Failed to register a consumer")
// 	  fmt.Println(6)
// 	  return msgs, err 
// }