package animesubs

import (
	"testing"
	"time"
)

func TestSubsInfoEqual(t *testing.T) {
	expected := true
	first := SubsInfo{
		Title:       "test",
		TimeUpdated: time.Unix(0, 0),
		Url:         "http://example.com",
	}
	second := SubsInfo{
		Title:       "test",
		TimeUpdated: time.Unix(0, 0),
		Url:         "http://example.com",
	}

	actual := first.Equal(second)
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
