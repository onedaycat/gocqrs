package cart

import "errors"

type Qty int

var (
	QtyNotBeZero = errors.New("Qty cannot be 0 or negative")
)

func (qty Qty) Increse(amount Qty) Qty {
	return qty + amount
}

func (qty Qty) Decrease(amount Qty) (Qty, error) {
	newqty := qty - amount
	if newqty < 1 {
		return 0, QtyNotBeZero
	}

	return newqty, nil
}
