package utils

import (
	"strings"
)

//SliceFind search string in slice of string
func SliceFind(slice []string, val string) (int, bool) {
	for i, item := range slice {

		if item == val {
			return i, true
		}
	}
	return -1, false
}

//SliceIndex Get position of element in slice
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

// DeleteEmpty ...
func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// FindChapters to find chapters from source url in the txt file
func FindChapters(cID string, cIDs []string, chapsText []string) (string, []string, []string) {
	ChapterIDs := []string{}
	var cTitle string
	for i := range cIDs {
		splitPath := strings.Split(cIDs[i], "/")
		index := SliceIndex(len(splitPath), func(i int) bool { return splitPath[i] == "chapters" })
		tChap := splitPath[index+1]
		ChapterIDs = append(ChapterIDs, tChap)

	}
	k, found := SliceFind(ChapterIDs, cID)

	if found {
		cTitle = chapsText[k]
	} else if !found && len(ChapterIDs) > 1 {
		cTitle = chapsText[0]
	} else if len(cID) == 0 {
		cTitle = "oneshot"
	}

	return cTitle, ChapterIDs, chapsText
}

func FindChaptersIDs(cIDs []string) []string {
	ChaptersIDs := []string{}
	for i := range cIDs {
		splitPath := strings.Split(cIDs[i], "/")
		index := SliceIndex(len(splitPath), func(i int) bool { return splitPath[i] == "chapters" })
		tChap := splitPath[index+1]
		ChaptersIDs = append(ChaptersIDs, tChap)
	}
	return ChaptersIDs
}
