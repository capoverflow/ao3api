package utils

import (
	"bufio"
	"log"
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
func ReadFile() []structs.ID {
	//var chaps []string
	var wID string
	var wChap string

	var extract []string
	var idsWork []structs.ID
	var idWork structs.ID
	//	var a int

	//archive := regexp.MustCompile(`https?:\/\/archiveofourown.org\/works\/[0-9]+`)
	archive := regexp.MustCompile(`https?://archiveofourown.org/works/[0-9]+(?:/chapters/[0-9]+)?`)
	lines, err := readLines("urls.txt")
	_ = err
	//log.Println("type", reflect.TypeOf(lines))
	for i := range lines {
		processedString := archive.FindString(lines[i])
		//log.Println(processedString)
		//log.Println(lines[i])
		extract = append(extract, processedString)
	}
	extract = DeleteEmpty(extract)
	//log.Println(len(extract))
	for i := range extract {
		//log.Println(i)
		//a = a+i
		splitPath := strings.Split(extract[i], "/")
		works := SliceIndex(len(splitPath), func(i int) bool { return splitPath[i] == "works" })
		chapters := SliceIndex(len(splitPath), func(i int) bool { return splitPath[i] == "chapters" })
		//log.Println(works, chapters)
		if len(splitPath) == 5 {
			//log.Println("Oneshot")
			//log.Println(splitPath[works+1])
			wID = splitPath[works+1]
			//log.Println(len(wID), wID)
			idWork.WorkID = wID
			idWork.ChapterID = ""
			idsWork = append(idsWork, idWork)
		} else if len(splitPath) == 7 {
			//print("MultiChapters")
			log.Println(splitPath[works+1], splitPath[chapters+1])
			wID = splitPath[works+1]
			wChap = splitPath[chapters+1]
			idWork.WorkID = wID
			idWork.ChapterID = wChap
			//log.Println(wID)
			idsWork = append(idsWork, idWork)
		}
		//log.Println(idsWork)
	}

	return idsWork
}

// FindChapters to find chapters from source url in the txt file
func FindChapters(cID string, cIDs []string, chapsText []string) (string, []string, []string) {
	//log.Println("... Debug FindChapters ...")
	// log.Println(cIDs)
	// log.Println(len(cIDs))
	ChapterIDs := []string{}
	var cTitle string
	for i := range cIDs {
		splitPath := strings.Split(cIDs[i], "/")
		index := SliceIndex(len(splitPath), func(i int) bool { return splitPath[i] == "chapters" })
		tChap := splitPath[index+1]
		ChapterIDs = append(ChapterIDs, tChap)
		// log.Printf("the n%d ChaptersID %s\n", i, ChapterIDs)

	}
	//log.Println("Chaps =", chaps)
	k, found := SliceFind(ChapterIDs, cID)
	//log.Println(len(chaps))

	if found {
		cTitle = chapsText[k]
		//log.Println(chapsText[k])
	} else if found == false && len(ChapterIDs) > 1 {
		cTitle = chapsText[0]
		//log.Println("Error")
	} else if len(cID) == 0 {
		cTitle = "oneshot"
		//log.Println("OneShot")
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
		// log.Printf("the n%d ChaptersID %s\n", i, ChapterIDs)

	}
	return ChaptersIDs
}
