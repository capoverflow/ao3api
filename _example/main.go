package main

import (
	"fmt"

	"gitlab.com/capoverflow/ao3api"
)

func main() {
	ParseWork := ao3api.ParseWork
	//adding test
	fmt.Println(ParseWork("17512811", "41254373"))
}
