package models

type FanficParams struct {
	Addr      string
	WorkID    string
	ChapterID string
	Debug     bool
	ProxyURLs []string
}

type Fanfic struct {
	URL             string      `json:"URL,omitempty"`
	WorkID          string      `json:"WorkID,omitempty"`
	ChapterID       string      `json:"ChapterID,omitempty"`
	ChapterTitle    string      `json:"ChapterTitle,omitempty"`
	Title           string      `json:"Title,omitempty"`
	Author          []string    `json:"Author,omitempty"`
	Published       string      `json:"Published,omitempty"`
	Updated         string      `json:"Updated,omitempty"`
	Words           string      `json:"Words,omitempty"`
	Chapters        string      `json:"Chapters,omitempty"`
	Comments        string      `json:"Comments,omitempty"`
	Kudos           string      `json:"Kudos,omitempty"`
	Bookmarks       string      `json:"Bookmarks,omitempty"`
	Hits            string      `json:"Hits,omitempty"`
	Fandom          []string    `json:"Fandom,omitempty"`
	Series          []string    `json:"Series,omitempty"`
	Summary         []string    `json:"Summary,omitempty"`
	ChaptersTitles  []string    `json:"ChaptersTitles,omitempty"`
	ChaptersIDs     []string    `json:"ChaptersIDs,omitempty"`
	Relationship    []string    `json:"Relationship,omitempty"`
	AlternativeTags []string    `json:"AlternativeTags,omitempty"`
	Freeform        []string    `json:"Freeform,omitempty"`
	Rating          []string    `json:"Rating,omitempty"`
	Character       []string    `json:"Character,omitempty"`
	Warning         []string    `json:"Warning,omitempty"`
	Category        []string    `json:"Category,omitempty"`
	Language        []string    `json:"Language,omitempty"`
	Additional      []string    `json:"Additional,omitempty"`
	Collection      []string    `json:"Collection,omitempty"`
	Challenge       []string    `json:"Challenge,omitempty"`
	Downloads       []Downloads `json:"Downloads,omitempty"`
	Status          string      `json:"Status,omitempty"`
}

type Downloads struct {
	FileType string
	Url      string
}
