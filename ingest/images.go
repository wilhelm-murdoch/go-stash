package ingest

type ImageJob struct {
	Path string
	Slug string
	Url  string
}

var ImageQueue chan ImageJob
