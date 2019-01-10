package models

//SyncLog stores info about sync log
type SyncLog struct {
	Model
	Content   string
	ProductID uint64 `gorm:"index:sync_product_idx"`
	Product   Product
}
