package main

import (
	"fmt"
	"log"

	"encoding/json"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println("%s: %s", msg, err)
	}
}

func amqpDial() (*amqp.Connection, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in amqpDial", r)
		}
	}()

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPassword, rabbitMQIp, rabbitMQPort)
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn, err
}

func Worker() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RabbitMQ Worker", r)
		}
	}()

	conn, err := amqpDial()
	defer conn.Close()
	if err != nil {
		return
	} else {
		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"DashboardEvents", // name
			true,              // durable
			false,             // delete when unused
			false,             // exclusive
			false,             // no-wait
			nil,               // arguments
		)
		failOnError(err, "Failed to declare a queue")
		err = ch.Qos(
			1,     // prefetch count
			0,     // prefetch size
			false, // global
		)
		failOnError(err, "Failed to set QoS")

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			false,  // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")

		forever := make(chan bool)

		go func() {
			fmt.Printf("closing: %s", <-conn.NotifyClose(make(chan *amqp.Error)))
			forever <- true
		}()

		go func() {
			fmt.Println("Subscribe to dashboard event")
			for d := range msgs {
				log.Printf("Received a dashboard event: %s", string(d.Body))
				d.Ack(false)

				var eventData EventData
				json.Unmarshal(d.Body, &eventData)

				go OnEvent(eventData.Tenent, eventData.Company, eventData.EventClass, eventData.EventType, eventData.EventCategory, eventData.SessionID, eventData.Parameter1, eventData.Parameter2, eventData.TimeStamp)

			}
			fmt.Println("Unsubscribe from dashboard event")
		}()

		log.Printf(" Dashboard Waiting for events. To exit press CTRL+C")
		<-forever
	}
}
