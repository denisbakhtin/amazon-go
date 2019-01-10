package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//Account stores db info about the user account
type Account struct {
	Model
	FirstName       string `form:"first_name" binding:"required"`
	LastName        string `form:"last_name" binding:"required"`
	Email           string `form:"email" binding:"required,email"`
	Password        string `form:"password"`
	PasswordConfirm string `form:"password_confirm" gorm:"-"`
	Role            string `form:"role" binding:"required"`
}

//BeforeCreate gorm hook
func (acc *Account) BeforeCreate() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(acc.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	acc.Password = string(hash)
	return nil
}

//PasswordIsValid checks a password against password policy
func PasswordIsValid(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("Password must be atleast 6 characters long")
	}
	return nil
}

//EncryptPassword applies creates bcrypt hash of password
func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
