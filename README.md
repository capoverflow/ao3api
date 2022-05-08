# A small golang api for archive of our own


## Sample program : 
```bash
cd ~/go/src
go mod init .
go mod edit -require github.com/capoverflow/ao3api@develop
go get -v -t ./...

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
```bash
go run .

1. The Contrary Corpse Salvage MuffinLance [50249441 50485448 50720243 50831173 51608758 52081273 54160033 54694417 60167875 61226317 61799353 63000286]
```
It return the chapters title, the title and the authors and the ids of all chapters.


## Roadmap: 

* Adding support for summary (already in the old api). WORKING (29-11-2020)
* Numbers of Kudos, Comments, Hits Working as of commit ()
* cli client (another project) [here](https://gitlab.com/capoverflow/ao3cmd)


