package models

//Company stores info about account company
type Company struct {
	Model
	Title       string
	Description string
	AccountID   uint64 `gorm:"index:company_account_idx"`
	Show        bool
	Account     Account
}
