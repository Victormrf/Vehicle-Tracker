package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventHub struct {
	routeService *RouteService
	mongoClient  *mongo.Client
	chDriverMoved chan *DriverMovedEvent
	freigthWriter *kafka.Writer
	simulatorWriter *kafka.Writer
}

func NewEventHub(routeService *RouteService, mongoClient *mongo.Client, chDriverMoved chan *DriverMovedEvent, freigthWriter *kafka.Writer, simulatorWriter *kafka.Writer) *EventHub {
	return &EventHub{
		routeService: routeService,
		mongoClient: mongoClient,
		chDriverMoved: chDriverMoved,
		freigthWriter: freigthWriter,
		simulatorWriter: simulatorWriter,
	}
}

func (eh *EventHub) HandleEvent(msg []byte) error {
	var basedEvent struct {
		EventName string `json:"event"`
	}
	err := json.Unmarshal(msg, &basedEvent)
	if err != nil {
		return fmt.Errorf("error unmarshalling event: %w", err)
	}

	switch basedEvent.EventName {
	case "RouteCreated":
		var event RouteCreatedEvent

		err := json.Unmarshal(msg, &event)
		if err != nil {
			return fmt.Errorf("error unmarshalling event: %w", err)
		}
		return eh.handleRouteCreated(event)

	case "DeliveryStarted":
		var event DeliveryStartedEvent
		err := json.Unmarshal(msg, &event)
		if err != nil {
			return fmt.Errorf("error unmarshalling event: %w", err)
		}
		return eh.handleDeliveryStarted(event)
	default: 
		return errors.New("Unknown event")
	}
}

func (eh *EventHub) handleRouteCreated(event RouteCreatedEvent) error {
	freigthCalculatedEvent, err := RouteCreatedHandler(&event, eh.routeService)
	if err != nil {
		return err
	}
	value, err := json.Marshal(freigthCalculatedEvent)
	if err != nil {
		return err
	}

	err = eh.freigthWriter.WriteMessages(context.Background(), kafka.Message{
		Key: []byte(freigthCalculatedEvent.RouteID),
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("error writing message: %w", err)
	}

	// publicar no apache kafka
	return nil
}

func (eh *EventHub) handleDeliveryStarted(event DeliveryStartedEvent) error {
	err := DeliveryStartedHandler(&event, eh.routeService, eh.chDriverMoved)
	if err != nil {
		return err
	}
	go eh.sendPositions() // goroute -- thread leve gerenciada pelo go

	// ler o canal e publicar no apache Kafka
	return nil
}

func (eh *EventHub) sendPositions() {
	for {
		select {
		case movedEvent := <- eh.chDriverMoved:
			value, err := json.Marshal(movedEvent)
			if err != nil {
				return
			}
			err = eh.simulatorWriter.WriteMessages(context.Background(), kafka.Message{
				Key: []byte(movedEvent.RouteID),
				Value: value,
			})
			if err != nil {
				return
			}
		case <- time.After(500 * time.Millisecond):
			return
		}
	}
}