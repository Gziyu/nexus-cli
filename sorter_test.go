package main

import (
	"testing"
	"time"
)

func Test_SortMixed(t *testing.T) {
	tags := []string{"latest", "1.0.1"}

	compareStringNumber := func(str1, str2 string) bool {
		return extractNumberFromString(str1) < extractNumberFromString(str2)
	}
	Compare(compareStringNumber).Sort(tags)

	if tags[0] != "1.0.1" && tags[1] != "latest" {
		t.Errorf("ordering incorrect when checking mixed tags")
	}
}

func Test_SortAllDigits(t *testing.T) {
	tags := []string{"1.2.1", "1.0.1"}

	compareStringNumber := func(str1, str2 string) bool {
		return extractNumberFromString(str1) < extractNumberFromString(str2)
	}
	Compare(compareStringNumber).Sort(tags)

	if tags[0] != "1.0.1" && tags[1] != "1.2.1" {
		t.Errorf("ordering incorrect in all digits tags")
	}
}

func Test_SortByTime(t *testing.T) {
	now := time.Now()
	tags := []TagInfo{
		{Name: "v1.0", Created: now.Add(-24 * time.Hour)},
		{Name: "v2.0", Created: now},
		{Name: "v1.5", Created: now.Add(-12 * time.Hour)},
	}

	// Test ascending sort (oldest first)
	compareTime := func(tag1, tag2 *TagInfo) bool {
		return tag1.Created.Before(tag2.Created)
	}
	TimeCompare(compareTime).Sort(tags)

	if tags[0].Name != "v1.0" || tags[1].Name != "v1.5" || tags[2].Name != "v2.0" {
		t.Errorf("time sorting incorrect: got %v, %v, %v", tags[0].Name, tags[1].Name, tags[2].Name)
	}
}
