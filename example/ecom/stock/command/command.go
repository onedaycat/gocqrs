package command

type CreateStockProduct struct {
	ProductID string
	Qty       int
}

type UpdateStockQuantity struct {
	ID  string
	Qty int
}

type RemoveStockProduct struct {
	ID string
}
