package helper

import (
	"math/rand"
	"time"
)

func ToPtrUint16(v uint16) *uint16 {
	return &v
}

func ToPtrInt16(v int16) *int16 {
	return &v
}

func RandomBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}
