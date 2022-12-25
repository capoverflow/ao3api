package utils

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func FindUrl(PageType string, urls []string) (max int) {
	var err error
	var nb []int
	var tnb int
	for _, url := range urls {
		slash := strings.HasSuffix(url, "/")
		if !slash {

			splitPath := strings.Split(url, "/")
			a := splitPath[len(splitPath)-1]
			re := regexp.MustCompile(fmt.Sprintf(`%s\?page=[0-9]+`, PageType))
			// // 			log.Println(a)
			z := re.FindAllString(a, -1)
			// log.Println(z)
			z = strings.Split(z[0], "=")
			// log.Println(z[1])
			tnb, err = strconv.Atoi(z[1])
			nb = append(nb, tnb)
			if err != nil {
				log.Println(err)
			}
		}
	}
	if len(nb) == 0 {
		return 0
	} else {
		return MaxIntSlice(nb)
	}
}

func MaxIntSlice(v []int) int {
	sort.Ints(v)
	return v[len(v)-1]
}

// SliceFind search string in slice of string
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

func RemoveDuplicates(strings []string) []string {
	// Create a map to store unique strings
	uniqueStrings := make(map[string]bool)

	// Iterate through the slice and add each string to the map
	for _, s := range strings {
		uniqueStrings[s] = true
	}

	// Create a slice to store the unique strings
	result := make([]string, 0, len(uniqueStrings))

	// Iterate through the map and add the keys (strings) to the slice
	for s := range uniqueStrings {
		result = append(result, s)
	}

	return result
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
