package vercontrol

import "time"

// Option 锁配置
type Option func(l *VerControl)

// SetExpire 设置锁过期时间
func SetExpire(t time.Duration) Option {
	return func(l *VerControl) {
		l.exp = t
	}
}
