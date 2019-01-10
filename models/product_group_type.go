package models

//ProductGroupType stores info about product group type
type ProductGroupType struct {
	Model
	CategoryID     uint64 `gorm:"index:pgt_category_idx"`
	ProductGroupID uint64 `gorm:"index:product_group_idx"`
	ProductTypeID  uint64 `gorm:"index:product_type_idx"`
	Category       Category
	ProductGroup   ProductGroup
	ProductType    ProductType
}
