package models

import (
	"time"
)

//TranslationRequest stores info about translation request
type TranslationRequest struct {
	Model
	Date    time.Time
	Count   int64
	Length  int64
	Allowed bool
}
