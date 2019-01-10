package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/denisbakhtin/amazon-go/models"
)

//AccountsGet shows all accounts
func AccountsGet(c *gin.Context) {
	var accounts []models.Account
	models.DB.Find(&accounts)

	H := DefaultH(c)
	H["Title"] = "Accounts"
	H["Accounts"] = accounts
	c.HTML(200, "admin/accounts/index", H)
}

//AccountsNewGet processes new account request
func AccountsNewGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	H := DefaultH(c)
	H["Title"] = "New account"
	H["Account"] = &models.Account{}
	H["Flash"] = flashes
	c.HTML(200, "admin/accounts/new", H)
}

//AccountsNewPost processes create account request
func AccountsNewPost(c *gin.Context) {
	account := models.Account{}
	if err := c.ShouldBind(&account); err != nil {
		sessionErrorAndRedirect(c, err, "/admin/new_account")
		return
	}
	if err := models.PasswordIsValid(account.Password); err != nil {
		sessionErrorAndRedirect(c, err, "/admin/new_account")
		return
	}
	if account.Password != account.PasswordConfirm {
		sessionErrorAndRedirect(c, fmt.Errorf("Password and password confirm do not match"), "/admin/new_account")
		return
	}

	if err := models.DB.Create(&account).Error; err != nil {
		sessionErrorAndRedirect(c, err, "/admin/new_account")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/accounts")
}

//AccountsEditGet processes edit account request
func AccountsEditGet(c *gin.Context) {
	id := c.Param("id")

	account := models.Account{}
	models.DB.First(&account, id)

	if account.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}
	account.Password = ""

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	H := DefaultH(c)
	H["Title"] = "Edit Account"
	H["Account"] = &account
	H["Flash"] = flashes
	c.HTML(200, "admin/accounts/edit", H)
}

//AccountsEditPost processes update account request
func AccountsEditPost(c *gin.Context) {
	id := c.Param("id")

	account := models.Account{}
	models.DB.First(&account, id)

	if account.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	model := models.Account{}
	if err := c.ShouldBind(&model); err != nil {
		fmt.Printf("%T: %+v", err, err)
		sessionErrorAndRedirect(c, err, fmt.Sprintf("/admin/edit_account/%d", account.ID))
		return
	}
	account.Password = "" //gorm will not update this field if it remains empty
	if len(model.Password) > 0 {
		//validate password
		if err := models.PasswordIsValid(model.Password); err != nil {
			sessionErrorAndRedirect(c, err, fmt.Sprintf("/admin/edit_account/%d", account.ID))
			return
		}
		if model.Password != model.PasswordConfirm {
			sessionErrorAndRedirect(c, fmt.Errorf("Password and password confirm do not match"), fmt.Sprintf("/admin/edit_account/%d", account.ID))
			return
		}
		if password, err := models.EncryptPassword(model.Password); err == nil {
			account.Password = password
		} else {
			sessionErrorAndRedirect(c, err, fmt.Sprintf("/admin/edit_account/%d", account.ID))
			return
		}
	}
	account.FirstName, account.LastName, account.Email, account.Role = model.FirstName, model.LastName, model.Email, model.Role

	//update only non-empty fields
	if err := models.DB.Model(&account).Updates(account).Error; err != nil {
		sessionErrorAndRedirect(c, err, fmt.Sprintf("/admin/edit_account/%d", account.ID))
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/accounts")
}

//AccountsDeletePost processes delete account request
func AccountsDeletePost(c *gin.Context) {
	id := c.Param("id")

	account := models.Account{}
	models.DB.First(&account, id)

	if account.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}
	if err := models.DB.Delete(&account).Error; err != nil {
		c.HTML(500, "errors/500", gin.H{"Error": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/accounts")
}
