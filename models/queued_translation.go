package models

//QueuedTranslation stores info about queued translations
type QueuedTranslation struct {
	Model
	ProductID  uint64 `gorm:"index:qt_product_idx"`
	LanguageID uint64 `gorm:"index:qt_language_idx"`
	Product    Product
	Language   Language
}
