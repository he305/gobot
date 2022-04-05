package stringutils

import (
	"math"

	levenshtein "github.com/ka-weihe/fast-levenshtein"
)

func GetLevenshteinDistance(first string, second string) int {
	return levenshtein.Distance(first, second)
}

func GetLevenshteinDistancePercent(first string, second string) int {
	distance := float64(GetLevenshteinDistance(first, second))
	biggerLength := math.Max(float64(len(first)), float64(len(second)))
	return int((biggerLength - distance) / biggerLength * 100)
}
