package response

import (
	"github.com/gogf/gf/v2/os/gtime"
)

var BootAt *gtime.Time

func init() {
	BootAt = gtime.Now()
}
func GetBootTime() int64 {
	if BootAt == nil {
		return -1
	}
	return gtime.Now().Timestamp() - BootAt.Timestamp()
}
func GetBootAt() string {
	return BootAt.String()
}
