package order

type Status string

const (
	PENDING  Status = "PENDING"
	PAID     Status = "PAID"
	REFUNED  Status = "REFUNED"
	REJECTED Status = "REJECTED"
)
