package util

import "strconv"

func StringToInt(val string, def int) int {
	result, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return result
}

func StringToInt64(val string) (int64, error) {
	result, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return int64(result), nil
}

func StringToUInt32(val string, def uint32) uint32 {
	result, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return uint32(result)
}
