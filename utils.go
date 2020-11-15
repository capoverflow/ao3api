package ao3api

import (
	"strings"
	//"fmt"
	//"regexp"
)

// DeleteEmpty delete empty String from slice
func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// SliceFind search string in slice of string used in main.go
func SliceFind(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// SliceIndex Get position of element in slice
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

//func oneChapter(splitPath []string) string {
//	index := SliceIndex(len(splitPath), func(i int) bool { return splitPath[i] == "works" })
//	tmp := splitPath[index-1]
//	return tmp
//}

// if works as multiple chapters return id of chapters in url
//func withChapters(splitPath []string) (string, string) {
//	index := SliceIndex(len(splitPath), func(i int) bool { return splitPath[i] == "chapters" })
//	tID := splitPath[index-1]
//	tChap := splitPath[index+1]
//	return tID, tChap
//}

// FindChapters to find chapters from source url in the txt file
func FindChapters(cID string, cIDs []string, chapsText []string) (string, []string) {
	//fmt.Println("... Debug ...")
	chaps := []string{}
	var cTitle string
	for i := range cIDs {
		splitPath := strings.Split(cIDs[i], "/")
		index := SliceIndex(len(splitPath), func(i int) bool { return splitPath[i] == "chapters" })
		tChap := splitPath[index+1]
		chaps = append(chaps, tChap)

	}
	//fmt.Println("Chaps =", chaps)
	k, found := SliceFind(chaps, cID)

	if found {
		cTitle = chapsText[k]
		//fmt.Println(chapsText[k])
	} else if !found && len(cID) != 0 {
		cTitle = "error"
		//fmt.Println("Error")
	}
	if len(cID) == 0 {
		cTitle = "oneshot"
		//fmt.Println("OneShot")
	}
	return cTitle, chaps
}
