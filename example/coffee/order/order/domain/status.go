package domain

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	OrderPlacedStatus    Status = "OrderPlaced"
	OrderAcceptedStatus  Status = "OrderAccepted"
	OrderStartedStatus   Status = "OrderStarted"
	OrderFinishedStatus  Status = "OrderFinished"
	OrderDeliveredStatus Status = "OrderDelivered"
)
