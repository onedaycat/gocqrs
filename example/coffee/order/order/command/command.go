package command

type PlaceOrder struct{}

type AcceptOrder struct {
	ID string `json:"id"`
}

type StartOrder struct {
	ID string `json:"id"`
}

type FinishOrder struct {
	ID string `json:"id"`
}

type DeliveryOrder struct {
	ID string `json:"id"`
}
