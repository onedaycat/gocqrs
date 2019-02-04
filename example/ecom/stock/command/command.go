package command

type CreateStockProduct struct {
	ProductID string
	Qty       int
}

type UpdateStockQuantity struct {
	ProductID string
	Qty       int
}

type RemoveStockProduct struct {
	ProductID string
}
