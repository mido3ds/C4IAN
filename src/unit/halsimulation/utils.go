package main

import (
	"math/rand"
)

func normal(mean, stdDev float64) float64 {
	return rand.NormFloat64()*stdDev + mean
}
