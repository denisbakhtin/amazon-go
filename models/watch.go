package models

//Watch stores info about product watches
type Watch struct {
	Model
	AccountID uint64 `gorm:"index:watch_account_idx"`
	ProductID uint64 `gorm:"index:watch_product_idx"`
	Account   Account
	Product   Product
}
