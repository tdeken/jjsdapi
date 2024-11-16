package progress

import "time"

// Option 进度配置
type Option func(l *Progress)

// SetExpire 设置进度条过期时间
func SetExpire(t time.Duration) Option {
	return func(l *Progress) {
		l.exp = t
	}
}
