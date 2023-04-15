package stringutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestAreSecondContainFirstCases(t *testing.T) {
	assert := assert.New(t)
	type testStruct struct {
		name string
		first string
		second string
		expected bool
	}

	testCases := []testStruct{
		{
			name: "Oshi No Ko #1",
			first: "[Oshi No Ko]",
			second: "edens zero - 28",
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := AreSecondContainsFirst(testCase.first, testCase.second)
			assert.Equal(testCase.expected, actual)
		})
	}
}

func TestIsStringContainsJapanese(t *testing.T) {
	assert := assert.New(t)
	type testStruct struct {
		name     string
		text     string
		expected bool
	}
	testCases := []testStruct{
		{
			name:     "Latin",
			text:     "some text",
			expected: false,
		},
		{
			name:     "Cyrillic",
			text:     "текст",
			expected: false,
		},
		{
			name:     "Only hiragana",
			text:     "おまえはもうしんでいる",
			expected: true,
		},
		{
			name:     "Only katakana",
			text:     "オマエハモウシンデイル",
			expected: true,
		},
		{
			name:     "Only kanji",
			text:     "具体的",
			expected: true,
		},
		{
			name:     "Mixed japanese",
			text:     "オマエはもう死んでいる",
			expected: true,
		},
		{
			name:     "Mixed with Latin",
			text:     "大切なもの protect my balls",
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := IsStringContainsJapanese(testCase.text)
			assert.Equal(testCase.expected, actual)
		})
	}

}
