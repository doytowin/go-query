package rdb

import (
	"fmt"
	. "github.com/doytowin/goooqo/core"
	"strconv"
	"time"
)

type TestEntity struct {
	Id         *int
	Username   *string
	Email      *string
	Mobile     *string
	CreateTime *time.Time
}

func (e TestEntity) GetTableName() string {
	return "t_user"
}

func (e TestEntity) GetId() any {
	return e.Id
}

func (e TestEntity) SetId(self any, id any) (err error) {
	v, ok := id.(int64)
	if !ok {
		s := fmt.Sprintf("%s", id)
		v, err = strconv.ParseInt(s, 10, 64)
	}
	if NoError(err) {
		self.(*TestEntity).Id = PInt(int(v))
	}
	return
}

type TestCond struct {
	Username *string
	Email    *string
	Mobile   *string
	TestAnd  *TestCond
}

type TestQuery struct {
	PageQuery
	EmailNull *bool
	TestOr    *TestCond
	Account   *string `condition:"(username = ? OR email = ?)"`
	Deleted   *bool
}
