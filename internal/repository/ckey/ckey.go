package ckey

const (
	Prefix = "jjsdapi:" //项目前缀

	UserPrefix = Prefix + "admin_user:" //用户缓存前缀
)

// ForLock 用来做临时锁
func ForLock(key string) string {
	return key + ":for_lock"
}
