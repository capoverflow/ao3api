# A small golang api for archive of our own


## Sample program : 
```bash
mkdir ao3api-example
cd ao3api-example
go mod init ao3api-example
go mod edit -require github.com/capoverflow/ao3api@master
```


```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	ao3 "github.com/capoverflow/ao3api"
	"github.com/capoverflow/ao3api/models"
)

func main() {

	fanfic, status, err := ao3.Fanfic(
		models.FanficParams{
			// Addr permit to use https://github.com/otwcode/otwarchive on your computer for developement or bug testing
			Addr:      "archiveofourown.org",
			WorkID:    "21116591",
			ChapterID: "50249441",
		},
	)
	if err != nil {
		log.Panic(err)
	}
	if status != 200 {
		log.Panicln("Site is down quitting")
	}

	b, err := json.MarshalIndent(fanfic, "", "  ")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(b))
}
``` 
```bash
ao3api-example ➤ go run .                                                                                                                                                                                     
ao3api: 2022/05/08 17:06:55 GetFirstChapterID.go:86: response received 200
{
  "URL": "http://archiveofourown.org/works/21116591",
  "WorkID": "21116591",
  "ChapterID": "50249441",
  "Title": "Salvage",
  "Author": [
    "MuffinLance"
  ],
  "Published": "2019-10-21",
  "Updated": "2021-05-27",
  "Words": "127176",
  "Chapters": "20/20",
  "Comments": "11744",
  "Kudos": "42261",
  "Bookmarks": "13993",
  "Hits": "958095",
  "Fandom": [
    "Avatar: The Last Airbender"
  ],
  "Summary": [
    "Mid-Season-One Zuko is held ransom by Chief Hakoda. Ozai's replies to the Water Tribe's demands are A+ Parenting. Hakoda is… deeply concerned, for this son that isn't his, and who might be safer among enemies than with his own father.",
    "Podfic and translations in French, German, Hungarian, Italian, Russian, and Spanish now available! See chapter one author notes for links."
  ],
  "Relationship": [
    "Hakoda \u0026 Zuko (Avatar)",
    "Zuko \u0026 The Southern Water Tribe",
    "Zuko \u0026 Responsible Adult Role Models"
  ],
  "AlternativeTags": [
    "Hakoda just wants to talk terms",
    "Ozai just wants a convenient barbarian to off his son in a politically expedient manner",
    "they are having a MINOR DISAGREEMENT on fatherhood",
    "Zuko is an Awkward Turtleduck",
    "who is very angry about being kept away from Avatar-hunting",
    "and also mildly concerned that someone is going to kill him in his sleep",
    "which is not stopping him from actively aggravating the enemy crew",
    "like a really growly puppy-kitten with a history of abuse",
    "let there be BONDING",
    "Hakuddles",
    "Hurt/Comfort",
    "Slowburn Adoption"
  ],
  "Downloads": [
    {
      "FileType": "azw3",
      "Url": "https://download.archiveofourown.org/downloads/21116591/Salvage.azw3?updated_at=1651610252"
    },
    {
      "FileType": "epub",
      "Url": "https://download.archiveofourown.org/downloads/21116591/Salvage.epub?updated_at=1651610252"
    },
    {
      "FileType": "mobi",
      "Url": "https://download.archiveofourown.org/downloads/21116591/Salvage.mobi?updated_at=1651610252"
    },
    {
      "FileType": "pdf",
      "Url": "https://download.archiveofourown.org/downloads/21116591/Salvage.pdf?updated_at=1651610252"
    },
    {
      "FileType": "html",
      "Url": "https://download.archiveofourown.org/downloads/21116591/Salvage.html?updated_at=1651610252"
    }
  ]
}
```
It return the fanfic struct in this example I transform it in json for a better display.


## Roadmap: 

* Adding support for summary (already in the old api). WORKING (29-11-2020)
* Numbers of Kudos, Comments, Hits Working as of commit ()
* cli client (another project) [here](https://gitlab.com/capoverflow/ao3cmd)


