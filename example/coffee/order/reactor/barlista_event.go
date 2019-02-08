package reactor

const (
	CoffeBrewStatusUpdatedEvent = "coffee.barlista.barlista.CoffeBrewStatusUpdated"
)

const (
	CoffeeBrewStartedStatus   = "CoffeeBrewStarted"
	CoffeeBrewFinishedStatus  = "CoffeeBrewFinished"
	CoffeeBrewDeliveredStatus = "CoffeeBrewDelivered"
)

type CoffeBrewStatusUpdated struct {
	Status string `json:"status"`
}
