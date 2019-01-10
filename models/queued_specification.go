package models

//QueuedSpecification stores info about queued specifications
type QueuedSpecification struct {
	Model
	ProductID uint64  `form:"product_id" binding:"required" gorm:"index:qs_product_idx"`
	Priority  uint64  `form:"priority"`
	Product   Product `form:"-" gorm:"association_autoupdate:false;association_autocreate:false"`
}
