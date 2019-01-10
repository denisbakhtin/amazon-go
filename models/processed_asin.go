package models

//ProcessedAsin stores info about processed asin
type ProcessedAsin struct {
	Model
	Asin string
	Log  string
}
