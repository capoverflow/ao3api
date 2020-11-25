# A small golang api for archive of our own


Sample program : 
```bash
cd ~/go/src
go mod init .
go mod edit -require gitlab.com/capoverflow/ao3api@develop
go get -v -t ./...

go run .
```


```go
package main

import (
	"fmt"

	ao3 "gitlab.com/capoverflow/ao3api"
)

func main() {
	Parser := ao3.ParseWorks
	fmt.Println(Parser("21116591", "50249441"))
}

``` 