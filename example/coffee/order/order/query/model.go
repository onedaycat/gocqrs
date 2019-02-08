package query

type Order struct {
	ID     string `json:"id" bson:"_id"`
	Price  int64  `json:"price" bson:"price"`
	Status string `json:"status" bson:"status"`
}
