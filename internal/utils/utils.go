package utils

import "strconv"

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

// StrToLongNumId 把字符串变回要处理的ID
func StrToLongNumId(v string) int64 {
	if v == "" {
		return 0
	}
	id, _ := strconv.ParseInt(v, 10, 64)
	return id
}

// LongNumIdToStr 把超长ID变成字符串，防止前端渲染大数字解析有问题
func LongNumIdToStr(v int64) string {
	if v == 0 {
		return ""
	}
	return strconv.FormatInt(v, 10)
}
