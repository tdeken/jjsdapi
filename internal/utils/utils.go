package utils

// Ternary 三元运算
func Ternary[T any](cond bool, t, f T) T {
	if cond {
		return t
	}
	return f
}

// PpDbLo page perPage 获取db limit和offset
func PpDbLo(page, perPage int32) (offset, limit int) {
	page = Ternary(page == 0, 1, page)
	perPage = Ternary(perPage == 0, 20, perPage)

	limit = int(perPage)
	offset = int((page - 1) * perPage)

	return
}
