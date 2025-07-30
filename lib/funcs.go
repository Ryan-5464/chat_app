package lib

import "strconv"

func ConvertInt64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func ConvertStringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
