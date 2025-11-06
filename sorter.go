package main

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Almost completely ripped off https://www.socketloop.com/tutorials/golang-natural-string-sorting-example

type Compare func(str1, str2 string) bool

type TagInfo struct {
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

func (cmp Compare) Sort(strs []string) {
	strSort := &strSorter{
		strs: strs,
		cmp:  cmp,
	}
	sort.Sort(strSort)
}

type strSorter struct {
	strs []string
	cmp  func(str1, str2 string) bool
}

func extractNumberFromString(str string) (num int) {
	strSlice := make([]string, 0)
	for _, v := range str {
		if unicode.IsDigit(v) {
			strSlice = append(strSlice, string(v))
		}
	}

	// If the tag was all non-digits, the strSlice would be empty (e.g., 'latest')
	// therefore just throw it to the end (1 << 32 is maxint)
	if len(strSlice) == 0 {
		return 1 << 32
	}

	num, err := strconv.Atoi(strings.Join(strSlice, ""))
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func (s *strSorter) Len() int { return len(s.strs) }

func (s *strSorter) Swap(i, j int) { s.strs[i], s.strs[j] = s.strs[j], s.strs[i] }

func (s *strSorter) Less(i, j int) bool { return s.cmp(s.strs[i], s.strs[j]) }

// TimeCompare compares tags based on their creation time
type TimeCompare func(tag1, tag2 *TagInfo) bool

func (cmp TimeCompare) Sort(tags []TagInfo) {
	tagSort := &tagSorter{
		tags: tags,
		cmp:  cmp,
	}
	sort.Sort(tagSort)
}

type tagSorter struct {
	tags []TagInfo
	cmp  func(tag1, tag2 *TagInfo) bool
}

func (s *tagSorter) Len() int { return len(s.tags) }

func (s *tagSorter) Swap(i, j int) { s.tags[i], s.tags[j] = s.tags[j], s.tags[i] }

func (s *tagSorter) Less(i, j int) bool { return s.cmp(&s.tags[i], &s.tags[j]) }
