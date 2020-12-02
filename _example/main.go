package main

import (
	"fmt"

	"gitlab.com/capoverflow/ao3api"
)

func main() {
	ParseWork := ao3api.ParseWork
	//adding test
	fmt.Println(ParseWork("23517349", "56400268"))
}
