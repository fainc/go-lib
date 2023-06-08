package model

import (
	"errors"
	"math"

	"github.com/gogf/gf/v2/database/gdb"
)

type helper struct {
	m *gdb.Model
}

func Helper(m *gdb.Model) *helper {
	if m == nil {
		panic("model error")
	}
	return &helper{m}
}

func (rec *helper) Create() (id int64, err error) {
	if id, err = rec.m.InsertAndGetId(); err != nil {
		return
	}
	return
}

func (rec *helper) Delete() (rows int64, err error) {
	ret, err := rec.m.Delete()
	if err != nil {
		return
	}
	rows, err = ret.RowsAffected()
	return
}

func (rec *helper) Update() (rows int64, err error) {
	ret, err := rec.m.Update()
	if err != nil {
		return
	}
	rows, err = ret.RowsAffected()
	return
}

func (rec *helper) Get(res interface{}) (err error) {
	err = rec.m.Scan(&res)
	return
}

func (rec *helper) MustGet(res interface{}) (err error) {
	err = rec.Get(res)
	if err != nil {
		return
	}
	if res == nil {
		err = errors.New("data does not exist")
	}
	return
}

func (rec *helper) CountWithPage(pageSize int) (total int, page int, err error) {
	total, err = rec.m.Count()
	if err != nil {
		return
	}
	if total <= 0 || pageSize <= 0 {
		return total, 1, nil
	}
	page = int(math.Ceil(float64(total) / float64(pageSize)))
	return
}
