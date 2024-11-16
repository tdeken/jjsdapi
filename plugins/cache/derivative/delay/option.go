package delay

type PushOption func(d *Delivery)

// UseTTL 生命周期
func UseTTL(ttl int) PushOption {
	return func(d *Delivery) {
		d.Ttl = ttl
	}
}

// UseAppId 应用ID
func UseAppId(id string) PushOption {
	return func(d *Delivery) {
		d.AppId = id
	}
}

// UseMsgId 生命周期
func UseMsgId(id string) PushOption {
	return func(d *Delivery) {
		d.MsgId = id
	}
}

// UseType 消息类型
func UseType(t string) PushOption {
	return func(d *Delivery) {
		d.Type = t
	}
}

type ConsumeOption func(c *Consumer)

// UserLogger 日志
func UserLogger(logger Logger) ConsumeOption {
	return func(c *Consumer) {
		c.logger = logger
	}
}
