package utils

import (
	"math/rand"
	"time"
)

type Utils struct {
}

func Number(min, max int) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	return float64(generator.Intn(max-min+1) + min)
}