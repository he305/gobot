package stringutils_test

import (
	"gobot/pkg/stringutils"
	"testing"
)

func TestGetLevenshteinDistanceCorrect(t *testing.T) {
	first := "kitten"
	second := "sitting"

	expected := 3

	actual := stringutils.GetLevenshteinDistance(first, second)

	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetLevenshteinDistancePercentCorrect(t *testing.T) {
	first := "shingeki no kyoujin"
	second := "shingeki+no+kyoujin"

	expected := 89

	actual := stringutils.GetLevenshteinDistancePercent(first, second)

	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
