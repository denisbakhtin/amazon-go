package models

//ProductTranslation stores info about product translation
type ProductTranslation struct {
	Model
	ProductID  uint64 `gorm:"index:pt_product_idx"`
	LanguageID uint64 `gorm:"index:pt_language_idx"`
	Title      string
	Content    string
	Product    Product
	Language   Language
}
