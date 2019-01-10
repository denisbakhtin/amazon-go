package models

import (
	"time"
)

//Feed stores info about rss feed
type Feed struct {
	Model
	URL            string
	Title          string
	Content        string
	Expires        time.Time
	Comments       string
	PagesNumber    int64
	PagesProcessed int64
	MorePages      bool
	State          bool
}

//GetURL returns feed url
func (f *Feed) GetURL() string {
	return "/" //not ready
}
