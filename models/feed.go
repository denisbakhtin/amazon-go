package models

import (
	"time"
)

//Feed stores info about rss feed
type Feed struct {
	Model
	URL      string
	Title    string
	Content  string
	Expires  time.Time
	Comments string
	State    bool
}

//GetURL returns feed url
func (f *Feed) GetURL() string {
	return "/" //not ready
}
