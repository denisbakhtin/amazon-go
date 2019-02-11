package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/denisbakhtin/amazon-go/models"
)

//PageGet processes GET /pages/1/slug
func PageGet(c *gin.Context) {
	id := c.Param("id")

	page := models.Page{}
	models.DB.First(&page, id)

	if page.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	//redirect to canonical url
	if c.Request.RequestURI != page.GetURL() {
		c.Redirect(http.StatusMovedPermanently, page.GetURL())
		return
	}

	H := DefaultH(c)
	H["Page"] = &page
	H["Title"] = page.Title
	c.HTML(200, "pages/show", H)
}

//PagesGet shows all pages
func PagesGet(c *gin.Context) {
	var pages []models.Page
	models.DB.Find(&pages)

	H := DefaultH(c)
	H["Title"] = "Pages"
	H["Pages"] = pages
	c.HTML(200, "admin/pages/index", H)
}

//PagesNewGet processes new page request
func PagesNewGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	page := models.Page{}
	H := DefaultH(c)
	H["Title"] = "New page"
	H["Page"] = &page
	H["Flash"] = flashes
	c.HTML(200, "admin/pages/new", H)
}

//PagesNewPost processes create page request
func PagesNewPost(c *gin.Context) {
	page := models.Page{}
	if err := c.ShouldBind(&page); err != nil {
		sessionErrorAndRedirect(c, err, "/admin/new_page")
		return
	}

	if err := models.DB.Create(&page).Error; err != nil {
		sessionErrorAndRedirect(c, err, "/admin/new_page")
		return
	}
	if page.Submit == submitAndViewTitle() {
		c.Redirect(http.StatusSeeOther, page.GetURL())
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/pages")
}

//PagesEditGet processes edit page request
func PagesEditGet(c *gin.Context) {
	id := c.Param("id")

	page := models.Page{}
	models.DB.First(&page, id)

	if page.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	H := DefaultH(c)
	H["Title"] = page.Title
	H["Page"] = &page
	H["Flash"] = flashes
	c.HTML(200, "admin/pages/edit", H)
}

//PagesEditPost processes update page request
func PagesEditPost(c *gin.Context) {
	id := c.Param("id")

	page := models.Page{}
	models.DB.First(&page, id)

	if page.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	vm := models.Page{}
	if err := c.ShouldBind(&vm); err != nil {
		sessionErrorAndRedirect(c, err, fmt.Sprintf("/admin/edit_page/%d", page.ID))
		return
	}
	page.Title, page.MetaKeywords, page.MetaDescription, page.Body, page.Show = vm.Title, vm.MetaKeywords, vm.MetaDescription, vm.Body, vm.Show

	if err := models.DB.Save(&page).Error; err != nil {
		sessionErrorAndRedirect(c, err, fmt.Sprintf("/admin/edit_page/%d", page.ID))
		return
	}
	if vm.Submit == submitAndViewTitle() {
		c.Redirect(http.StatusSeeOther, page.GetURL())
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/pages")
}

//PagesDeletePost processes delete page request
func PagesDeletePost(c *gin.Context) {
	id := c.Param("id")

	page := models.Page{}
	models.DB.First(&page, id)

	if page.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}
	if err := models.DB.Delete(&page).Error; err != nil {
		c.HTML(500, "errors/500", gin.H{"Error": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/pages")
}
