package helper

import "math"

func PageOffset(limit int, page int) (offset int) {
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}
	return
}

func PageTotal(limit int, page int, count int64) (totalPage int, pageIsEnd int) {
	if count > 0 {
		totalPage = int(math.Ceil(float64(count) / float64(limit)))
	}
	if page >= totalPage {
		pageIsEnd = 1
	}
	return
}