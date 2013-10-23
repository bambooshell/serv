package util

import (
	"fmt"
	"time"
)

//represent time to string
//in the form of yy-mm-dd HH-MM-SS
func Time2Str(t *time.Time) (str string) {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	str = fmt.Sprintf("%4d-%02d-%02d %02d-%02d-%02d", year, month, day, hour, min, sec)

	return str
}
