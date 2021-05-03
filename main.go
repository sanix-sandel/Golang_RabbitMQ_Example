package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err, "Dialing failed to RabbitMQ broker")
	defer conn.Close()

	channel, err := conn.Channel()
	handleError(err, "Fetching channel failed")
	defer channel.Close()

	//creating a queue
	testQueue, err := channel.QueueDeclare(
		"test", //name of the queue
		false,  //Message is persisted or not
		false,  //Delete message when unused
		false,  //Exclusive
		false,  //No waiting time
		nil,    //Extra args
	)
	handleError(err, "Queue creation failed")

	//preparing amqp message
	serverTime := time.Now()
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(serverTime.String()),
	}

	//publishing the message to the queue
	err = channel.Publish(
		"",             //exchange
		testQueue.Name, //routing key (Queue)
		false,          //mandatory
		false,          //immediate
		message,        //immediate
	)
	handleError(err, "Failed to publish a message")
	log.Println("Successfully published a message to the queue")
}
