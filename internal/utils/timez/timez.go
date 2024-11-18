package timez

import "time"

// TableDateTime 后台列表时间
func TableDateTime(t int64) string {
	return time.Unix(t, 0).Format(time.DateTime)
}
