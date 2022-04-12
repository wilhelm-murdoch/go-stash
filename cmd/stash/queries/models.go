package queries

// Author
type Author struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Tagline     string `json:"tag_line"`
	Photo       string `json:"photo"`
	CoverImage  string `json:"cover_image"`
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
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Tagline string `json:"tag_line"`
}

// Post
type Post struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	CUID        string `json:"cuid"`
	Type        string `json:"type"`
	DateAdded   string `json:"date_added"`
	DateUpdated string `json:"date_updated"`
	Content     string `json:"content"`
	Brief       string `json:"brief"`
	CoverImage  string `json:"cover_image"`
	Author      Author
	Tags        []Tag
}
