package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//SignInGet handles GET /signin route
func SignInGet(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "Sign in"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "auth/signin", h)
}

//SignInPost processes POST /signin
func SignInPost(c *gin.Context) {
	signin := models.SignIn{}
	session := sessions.Default(c)
	if err := c.ShouldBind(&signin); err != nil {
		sessionErrorAndRedirect(c, err, "/signin")
		log.Printf("Wrong login or password, IP: %s", c.Request.RemoteAddr)
		return
	}

	account := models.Account{}
	models.DB.Where("email = lower(?)", signin.Email).First(&account)

	if account.ID == 0 {
		sessionErrorAndRedirect(c, fmt.Errorf("Wrong login or password"), "/signin")
		log.Printf("Wrong login or password, IP: %s", c.Request.RemoteAddr)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(signin.Password)); err != nil {
		sessionErrorAndRedirect(c, fmt.Errorf("Wrong login or password"), "/signin")
		log.Printf("Wrong login or password, IP: %s", c.Request.RemoteAddr)
		return
	}

	session.Set(accountIDKey, account.ID)
	session.Set(accountRoleKey, account.Role)
	session.Save()
	c.Redirect(http.StatusFound, "/admin")
}

//SignUpGet handles GET /signup route
func SignUpGet(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "Sign up"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "auth/signup", h)
}

//SignUpPost processes POST /signup
func SignUpPost(c *gin.Context) {
	signup := models.SignUp{}
	session := sessions.Default(c)
	if err := c.ShouldBind(&signup); err != nil {
		sessionErrorAndRedirect(c, err, "/signup")
		return
	}
	if err := models.PasswordIsValid(signup.Password); err != nil {
		sessionErrorAndRedirect(c, err, "/signup")
		return
	}
	if signup.Password != signup.PasswordConfirmation {
		sessionErrorAndRedirect(c, fmt.Errorf("Password and password confirm do not match"), "/signup")
		return
	}

	account := models.Account{}
	models.DB.Where("email = lower(?)", signup.Email).First(&account)

	if account.ID != 0 {
		sessionErrorAndRedirect(c, fmt.Errorf("This email has already been taken"), "/signup")
		return
	}

	account = models.Account{Email: signup.Email, Password: signup.Password, Role: config.UserRole}
	if err := models.DB.Create(&account).Error; err != nil {
		sessionErrorAndRedirect(c, err, "/signup")
		log.Println(err.Error())
		return
	}

	session.Set(accountIDKey, account.ID)
	session.Set(accountRoleKey, account.Role)
	session.Save()
	c.Redirect(http.StatusSeeOther, "/")
}

//SignOutPost processes POST /signout
func SignOutPost(c *gin.Context) {
	session := sessions.Default(c)

	session.Delete(accountIDKey)
	session.Delete(accountRoleKey)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/")
}
