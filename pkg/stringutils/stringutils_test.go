package stringutils

import (
	"testing"
)

func TestGetLevenshteinDistanceCorrect(t *testing.T) {
	first := "kitten"
	second := "sitting"

	expected := 3

	actual := GetLevenshteinDistance(first, second)

	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetLevenshteinDistancePercentCorrect(t *testing.T) {
	first := "shingeki no kyoujin"
	second := "shingeki+no+kyoujin"

	expected := 89

	actual := GetLevenshteinDistancePercent(first, second)

	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestAreSecondContainsFirstTrue(t *testing.T) {
	first := "shinGekI nO kyoujin"
	second := "shingeki no kyoujin - final season 04"

	expected := true
	actual := AreSecondContainsFirst(first, second)

	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestAreSecondContainsFirstFalse(t *testing.T) {
	first := "Jojo Golden wind"
	second := "Jojo no Kimyou na Bouken"

	expected := false
	actual := AreSecondContainsFirst(first, second)

	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
