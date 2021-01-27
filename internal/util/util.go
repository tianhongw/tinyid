package util

import (
	"math/rand"
	"time"
)

func RandomInt(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(n)
}
