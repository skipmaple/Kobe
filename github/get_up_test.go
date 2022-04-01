package github

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	year, month, day := time.Now().In(loc).Date()
	since := time.Date(year, month, day, 0, 0, 0, 0, loc)
	fmt.Printf("%v\n", since)
}
