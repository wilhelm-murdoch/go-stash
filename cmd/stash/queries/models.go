package queries

// Author
type Author struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Tagline     string `json:"tagLine"`
	Photo       string `json:"photo"`
	CoverImage  string `json:"coverImage"`
	SocialMedia SocialMedia
}

// SocialMedia
type SocialMedia struct {
	Twitter       string `json:"twitter"`
	Github        string `json:"github"`
	StackOverflow string `json:"stackoverflow"`
	LinkedIn      string `json:"linkedin"`
	Google        string `json:"google"`
	Website       string `json:"website"`
	Facebook      string `json:"facebook"`
}

// Tag
type Tag struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	Logo string `json:"logo"`
}

// Post
type Post struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	DateAdded   string `json:"dateAdded"`
	DateUpdated string `json:"dateUpdated"`
	CUID        string `json:"cuid"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Brief       string `json:"brief"`
	CoverImage  string `json:"coverImage"`
	Author      Author
	Tags        []Tag
}
