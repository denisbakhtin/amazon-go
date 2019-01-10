package models

//Specification stores info about product specification
type Specification struct {
	Model
	ProductID   uint64 `gorm:"index:spec_product_idx"`
	Description string
	Feature     string
	Product     Product
}
