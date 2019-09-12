package order

import (
	"math/rand"
	"time"
)

func GenerateNumber() string {
	timeNow := time.Now().Format("20060102150405")
	tempRand := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k",
		"l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v",
		"w", "x", "y", "z",
	}
	number := "W" + timeNow
	for i := 0; i <= 3; i++ {
		number = number + tempRand[rand.Intn(26)]
	}
	return number
}
