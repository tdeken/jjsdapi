package bucket

import "time"

type LastingOption func(lb *LastingBucket)

// UseCheckValue 产生的一个值是否可以入桶
func UseCheckValue(check CheckValue) LastingOption {
	return func(lb *LastingBucket) {
		lb.check = append(lb.check, check)
	}
}

// UserLastingLen 设置桶的最大长度
func UserLastingLen(bl int64) LastingOption {
	return func(lb *LastingBucket) {
		lb.bl = bl
	}
}

// UserLastingThreshold 设置桶的生成阈值
//
//	当剩余长度小于或等于这个值得时候才填充桶，必须大于0，小于或等于0的时候无效
func UserLastingThreshold(trd int64) LastingOption {
	return func(lb *LastingBucket) {
		lb.threshold = trd
	}
}

// UserLife 设置桶的静止生命周期
//
//	每次设置获取完值之后，桶如果在这个时间内未被使用，将直接销毁
func UserLife(life time.Duration) LastingOption {
	return func(lb *LastingBucket) {
		lb.life = life
	}
}
