package models

//Department stores info about product department
type Department struct {
	Model
	Title string
	Count int `gorm:"-"` //product count
}
