package stringutils

import (
	"math"
	"regexp"
	"strings"
)

func GetLevenshteinDistance(first string, second string) int {
	return levenshtein([]rune(first), []rune(second))
}

func GetLevenshteinDistancePercent(first string, second string) int {
	distance := float64(GetLevenshteinDistance(first, second))
	biggerLength := math.Max(float64(len(first)), float64(len(second)))
	return int((biggerLength - distance) / biggerLength * 100)
}

func AreSecondContainsFirst(first string, second string) bool {
	first = strings.ToLower(first)
	second = strings.ToLower(second)

	matched, _ := regexp.MatchString(first, second)
	return matched
}

func LowerAndTrimText(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return text
}

func IsStringContainsJapanese(text string) bool {
	byteText := []byte(text)
	var regexHiragana = regexp.MustCompile(`[^\P{Hiragana}\p{N}]`)

	if regexHiragana.Match(byteText) {
		return true
	}

	var regexKanji = regexp.MustCompile(`[^\P{Han}\p{N}]`)
	if regexKanji.Match(byteText) {
		return true
	}

	var regexKatakana = regexp.MustCompile(`[^\P{Katakana}\p{N}]`)
	return regexKatakana.Match(byteText)
}
