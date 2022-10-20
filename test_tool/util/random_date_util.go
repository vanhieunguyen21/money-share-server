package util

import (
	"math/rand"
	"time"
)

func RandomDateTillNow(fromYear int) time.Time {
	min := time.Date(fromYear, time.January, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func RandomDateBetweenYears(fromYear int, toYear int) time.Time {
	if toYear < fromYear {
		panic("toYear must not be smaller than fromYear")
	}
	min := time.Date(fromYear, time.January, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(toYear, time.January, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
