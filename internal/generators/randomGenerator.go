package generators

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomInt32(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
