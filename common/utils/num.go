package utils

import "strconv"

func InterfaceToInt64(value interface{}) int64 {
	if value == nil {
		return 0
	}
	switch value.(type) {
	case int:
		return int64(value.(int))
	case int8:
		return int64(value.(int8))
	case int16:
		return int64(value.(int16))
	case int32:
		return int64(value.(int32))
	case int64:
		return value.(int64)
	case uint:
		return int64(value.(uint))
	case uint8:
		return int64(value.(uint8))
	case uint16:
		return int64(value.(uint16))
	case uint32:
		return int64(value.(uint32))
	case uint64:
		return int64(value.(uint64))
	case float32:
		return int64(value.(float32))
	case float64:
		return int64(value.(float64))
	case string:
		i, _ := strconv.ParseInt(value.(string), 10, 64)
		return i
	default:
		return 0
	}
}
