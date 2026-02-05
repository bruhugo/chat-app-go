package utils

import (
	"strconv"
)

func ConvertAnyToInt64(a any) (int64, bool) {
	switch t := a.(type) {
	case string:
		r, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0, false
		}

		return r, true

	case int:
		return int64(t), true
	case int64:
		return t, true
	default:
		return 0, false
	}
}
