package ao3structs

type Work struct {
	URL             string   `json:"URL,omitempty"`
	WorkID          string   `json:"WorkID,omitempty"`
	ChapterID       string   `json:"ChapterID,omitempty"`
	ChapterTitle    string   `json:"ChapterTitle,omitempty"`
	Title           string   `json:"Title,omitempty"`
	Author          string   `json:"Author,omitempty"`
	Published       string   `json:"Published,omitempty"`
	Updated         string   `json:"Updated,omitempty"`
	Words           string   `json:"Words,omitempty"`
	Chapters        string   `json:"Chapters,omitempty"`
	Comments        string   `json:"Comments,omitempty"`
	Kudos           string   `json:"Kudos,omitempty"`
	Bookmarks       string   `json:"Bookmarks,omitempty"`
	Hits            string   `json:"Hits,omitempty"`
	Fandom          string   `json:"Fandom,omitempty"`
	Summary         []string `json:"Summary,omitempty"`
	ChaptersTitles  []string `json:"ChaptersTitles,omitempty"`
	ChaptersIDs     []string `json:"ChaptersIDs,omitempty"`
	Relationship    []string `json:"Relationship,omitempty"`
	AlternativeTags []string `json:"AlternativeTags,omitempty"`
}

type ID struct {
	WorkID    string
	ChapterID string
}

type Fanfic struct {
	WorkID    string
	ChapterID string
	Debug     bool
}
