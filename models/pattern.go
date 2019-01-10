package models

//Pattern stores info about category pattern
type Pattern struct {
	Model
	CategoryID     uint64 `gorm:"index:pattern_category_idx"`
	Title          string
	TitlePattern   string
	BrowsePattern  string
	BrowseAndTitle bool
	Priority       float64
	Category       Category
}
