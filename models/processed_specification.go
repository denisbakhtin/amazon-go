package models

//ProcessedSpecification stores info about processed product specification
type ProcessedSpecification struct {
	Model
	ProductID uint64  `form:"product_id" gorm:"index:ps_product_idx"`
	Product   Product `form:"-" gorm:"association_autoupdate:false;association_autocreate:false"`
}
