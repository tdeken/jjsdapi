package lock

import "time"

// Option 锁配置
type Option func(l *ILock)

// SetExpire 设置锁过期时间
func SetExpire(t time.Duration) Option {
	return func(l *ILock) {
		l.exp = t
	}
}

// SetValue 自定义锁值
func SetValue(val string) Option {
	return func(l *ILock) {
		l.value = val
	}
}
