package response

import (
	"github.com/gogf/gf/v2/os/gtime"
)

var BootAt *gtime.Time

func BootMark() {
	BootAt = gtime.Now()
}
func GetBootTime() int64 {
	return gtime.Now().Timestamp() - BootAt.Timestamp()
}
func GetBootAt() string {
	return BootAt.String()
}
