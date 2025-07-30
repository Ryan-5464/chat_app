package lib

import "strconv"

func ConvertInt64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
