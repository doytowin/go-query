package mongodb

import (
	"github.com/doytowin/goquery"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryEntity struct {
	Id   primitive.ObjectID `bson:"_id,omitempty"`
	Item string
	Size struct {
		H   float64
		W   float64
		Uom string
	}
	Qty    int
	Status string
}

func (r InventoryEntity) Database() string {
	return "doytowin"
}

func (r InventoryEntity) Collection() string {
	return "inventory"
}

type InventoryQuery struct {
	goquery.PageQuery
}