package utils

import "errors"

// CalculateOffset 计算分页偏移量
func CalculateOffset(page, pageSize int) (int, error) {

	if page < 1 || pageSize < 1 || pageSize > 100 {
		return 0, errors.New("分页参数不合法")
	}

	maxInt := int(^uint(0) >> 1)

	if page-1 > maxInt/pageSize {
		return 0, errors.New("分页偏移量溢出")
	}

	offset := (page - 1) * pageSize

	return offset, nil
}
