package main

import (
	"github.com/doytowin/goooqo"
	"github.com/doytowin/goooqo/mongodb"
	. "go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryQuery struct {
	goooqo.PageQuery
	QtyGt *int
}

type SizeDoc struct {
	H   float64 `json:"h,omitempty" bson:"h"`
	W   float64 `json:"w,omitempty" bson:"w"`
	Uom string  `json:"uom,omitempty" bson:"uom"`
}

type InventoryEntity struct {
	mongodb.MongoId `bson:",inline"`
	Item            string  `json:"item,omitempty" bson:"item"`
	Size            SizeDoc `json:"size" bson:"size"`
	Qty             int     `json:"qty,omitempty" bson:"qty"`
	Status          string  `json:"status,omitempty" bson:"status"`
}

func (r InventoryEntity) Database() string {
	return "doytowin"
}

func (r InventoryEntity) Collection() string {
	return "inventory"
}

func (q InventoryQuery) BuildFilter() []D {
	d := make([]D, 0, 10)
	if q.QtyGt != nil {
		d = append(d, D{{"qty", D{{"$gt", q.QtyGt}}}})
	}
	return d
}