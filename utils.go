package neoao3api

import (
	"bufio"
	"fmt"
	"os"

	//	"reflect"
	"regexp"
	"strings"
	//	"log"
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

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// ReadFile read from file
func ReadFile() []id {
	//var chaps []string
	var wID string
	var wChap string

	var extract []string
	var idsWork []id
	var idWork id
	//	var a int

	//archive := regexp.MustCompile(`https?:\/\/archiveofourown.org\/works\/[0-9]+`)
	archive := regexp.MustCompile(`https?://archiveofourown.org/works/[0-9]+(?:/chapters/[0-9]+)?`)
	lines, err := readLines("urls.txt")
	_ = err
	//fmt.Println("type", reflect.TypeOf(lines))
	for i := range lines {
		processedString := archive.FindString(lines[i])
		//fmt.Println(processedString)
		//fmt.Println(lines[i])
		extract = append(extract, processedString)
	}
	extract = DeleteEmpty(extract)
	//fmt.Println(len(extract))
	for i := range extract {
		//fmt.Println(i)
		//a = a+i
		splitPath := strings.Split(extract[i], "/")
		works := SliceIndex(len(splitPath), func(i int) bool { return splitPath[i] == "works" })
		chapters := SliceIndex(len(splitPath), func(i int) bool { return splitPath[i] == "chapters" })
		//fmt.Println(works, chapters)
		if len(splitPath) == 5 {
			//fmt.Println("Oneshot")
			//fmt.Println(splitPath[works+1])
			wID = splitPath[works+1]
			//fmt.Println(len(wID), wID)
			idWork.WorkID = wID
			idWork.ChapterID = ""
			idsWork = append(idsWork, idWork)
		} else if len(splitPath) == 7 {
			//print("MultiChapters")
			fmt.Println(splitPath[works+1], splitPath[chapters+1])
			wID = splitPath[works+1]
			wChap = splitPath[chapters+1]
			idWork.WorkID = wID
			idWork.ChapterID = wChap
			//fmt.Println(wID)
			idsWork = append(idsWork, idWork)
		}
		//fmt.Println(idsWork)
	}

	return idsWork
}

// FindChapters to find chapters from source url in the txt file
func FindChapters(cID string, cIDs []string, chapsText []string) (string, []string) {
	//fmt.Println("... Debug FindChapters ...")
	//fmt.Println(cIDs)
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
	//fmt.Println(len(chaps))

	if found {
		cTitle = chapsText[k]
		//fmt.Println(chapsText[k])
	} else if found == false && len(chaps) > 1 {
		cTitle = chapsText[0]
		fmt.Println("Error")
	} else if len(cID) == 0 {
		cTitle = "oneshot"
		//fmt.Println("OneShot")
	}

	return cTitle, chaps
}
