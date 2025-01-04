package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Victormrf/Vehicle-tracker/simulator/internal"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoStr := "mongodb://admin:admin@localhost:27017/routes?authSource=admin"
	mongoConnection, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoStr))
	if err != nil {
		panic(err)
	}

	freigthService := internal.NewFreigthService()
	routeService := internal.NewRouteService(mongoConnection, freigthService)

	chDriverMoved := make(chan *internal.DriverMovedEvent)
	kafkaBroker := "localhost:9092"

	freightWriter := &kafka.Writer{
		Addr: kafka.TCP(kafkaBroker),
		Topic: "freight",
		Balancer: &kafka.LeastBytes{},
	}

	simulatorWriter := &kafka.Writer{
		Addr: kafka.TCP(kafkaBroker),
		Topic: "simulator",
		Balancer: &kafka.LeastBytes{},
	}

	routeReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    "route",
		GroupID:  "route-group",
	})

	hub := internal.NewEventHub(routeService, mongoConnection, chDriverMoved, freightWriter, simulatorWriter)

	fmt.Println("Starting simulator")
	for {
		m, err := routeReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}

		go func(msg []byte) {
			err = hub.HandleEvent(m.Value)
			if err != nil {
				log.Printf("error handling event: %v", err)
			}
		}(m.Value)
	}
}

