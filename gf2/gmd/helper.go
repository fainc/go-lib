package gmd

//  gf2 model helper
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
func (rec *helper) MustGet(res interface{}) (err error) {
	err = rec.m.Scan(&res)
	if err != nil {
		return
	}
	if res == nil {
		return errors.New("the query data does not exist")
	}
	return
}
