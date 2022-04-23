package models

type FanficParams struct {
	Addr      string
	WorkID    string
	ChapterID string
	Debug     bool
	ProxyURLs []string
}

//Work ..
type Work struct {
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
	Summary         []string    `json:"Summary,omitempty"`
	ChaptersTitles  []string    `json:"ChaptersTitles,omitempty"`
	ChaptersIDs     []string    `json:"ChaptersIDs,omitempty"`
	Relationship    []string    `json:"Relationship,omitempty"`
	AlternativeTags []string    `json:"AlternativeTags,omitempty"`
	Downloads       []Downloads `json:"Downloads,omitempty"`
}

type Downloads struct {
	FileType string
	Url      string
}

//FanficID
type FanficID struct {
	WorkID    string
	ChapterID string
}

type Fanfic struct {
	WorkID    string
	ChapterID string
	Debug     bool
}

type Search struct {
	AnyField         string
	Title            string
	Author           string
	Oneshot          bool
	Language         string
	CompletionStatus bool
	Fandoms          string
	Relationship     string
}

type UserParams struct {
	Addr      string
	Username  string
	Debug     bool
	ProxyURLs []string
}

type User struct {
	Username  string
	Profile   UserProfile
	Works     []Work
	Bookmarks []Work
	Gift      []Work
}

type UserProfile struct {
	Pseuds   string
	JoinDate string
	Email    string
	Bio      string
}
