package internal

type RouteCreatedEevent struct {
	EventName  string       `json:"event"`
	RouteID    string       `json:"id"`
	Distance   int          `json:"distance"`
	Directions []Directions `json:"directions"`
}

func NewRouteCreatedEvent(routeID string, distance int, directions []Directions) *RouteCreatedEevent {
	return &RouteCreatedEevent{
		EventName:  "routeCreated",
		RouteID:    routeID,
		Distance:   distance,
		Directions: directions,
	}
}

type FreightCalculatedEvent struct {
	EventName string  `json:"event"`
	RouteID   string  `json:"id"`
	Amount    float64 `json:"amount"`
}

func NewFreightCalculatedEvent(routeID string, amount float64) *FreightCalculatedEvent {
	return &FreightCalculatedEvent{
		EventName: "freightCalculated",
		RouteID:   routeID,
		Amount:    amount,
	}
}

type DeliveryStartedEvent struct {
	EventName string `json:"event"`
	RouteID   string `json:"route_id"`
} 

func NewDeliveryStartedEvent(routeID string) *DeliveryStartedEvent {
	return &DeliveryStartedEvent{
		EventName: "DeliveryStarted",
		RouteID: routeID,
	}
}

func RouteCreatedHandler(event *RouteCreatedEevent, routeService *RouteService) (*FreightCalculatedEvent, error) {
	route := NewRoute(event.RouteID, event.Distance, event.Directions)
	routeCreated, err := routeService.CreateRoute(route)
	if err != nil {
		return nil, err
	}
	freightCalculatedEvent := NewFreightCalculatedEvent(routeCreated.ID, routeCreated.FreightPrice)
	return freightCalculatedEvent, nil
}
