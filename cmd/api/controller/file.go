package controller

import (
	_ "bucket_file/constant"
	_ "fmt"
	"log"
	_ "net/http"
	_ "bucket_file/server"
	_ "github.com/streadway/amqp"
	_ "github.com/gin-gonic/gin"
)

type FileForm struct {
	FileName string `json:"filename" bson:"filename"`
}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Fatalf("%s: %s", msg, err)
	}
}

// func Create_file(c *gin.Context) {
// 	conn, err := amqp.Dial("amqp://test:test@192.168.30.94:5672/")
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 	"upload file", // name
// 	false,   	 // durable
// 	false,   	 // delete when usused
// 	false,   	 // exclusive
// 	false,   	 // no-wait
// 	nil,    	 // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")

// 	msgs, err := ch.Consume(
// 		q.Name, // queue
// 		"",     // consumer
// 		true,   // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // args
// 	  )
// 	failOnError(err, "Failed to register a consumer")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"msg": constant.SERVER_ERROR_500,
// 		})
// 		return
// 	}
// 	// fmt.Println(msgs)
// 	stopChan := make(chan bool)
// 	go func() {
// 		for d := range msgs {
// 			log.Printf("Received a message: %s", d.Body)
// 			fmt.Printf("Received a message: %s", d.Body)
// 		}
// 	}()

// 	<-stopChan
// }
	// userID := c.MustGet("userID").(string)
	// user, _, _ := container.MongoClient.User().FindById(userID)
	// res, code := container.MongoClient.File().FileRequest("PUT", user.Username, postData.FileName)
	// if code != 200 {
	// 	errorResponse := container.MongoClient.Bucket().ErrorXml2Json(res)
	// 	c.JSON(code, gin.H{
	// 		"msg":    "error",
	// 		"result": errorResponse,
	// 	})
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "Created file",
	// })
// }

// func GetFile(c *gin.Context) {
// 	fileName := c.Query("filename")
// 	userID := c.MustGet("userID").(string)
// 	user, _, _ := container.MongoClient.User().FindById(userID)
// 	res, code := container.MongoClient.File().FileRequest("GET", user.Username, fileName)
// 	fmt.Println(res)
// 	if code != 200 {
// 		errorResponse := container.MongoClient.Bucket().ErrorXml2Json(res)
// 		c.JSON(code, gin.H{
// 			"msg":    "error",
// 			"result": errorResponse,
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"msg": "file",
// 	})
// }
