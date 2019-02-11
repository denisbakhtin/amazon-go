package controllers

import (
	"fmt"
	"math"
	"net/http"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/denisbakhtin/amazon-go/utility"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//CategoryGet processes GET /categories/1/slug
func CategoryGet(c *gin.Context) {
	id := c.Param("id")

	category := models.Category{}
	models.DB.First(&category, id)

	if category.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	//redirect to canonical url
	if c.Request.URL.Path != category.GetURL() {
		c.Redirect(http.StatusMovedPermanently, category.GetURL())
		return
	}

	category.LoadChildren()
	categoryIDs := category.GetCategoryIDs()

	totalCount := 0
	dbQuery := models.DB.Model(models.Product{}).Where("category_id IN(?) AND available = true", categoryIDs)
	dbQuery = applyProductFilters(c, dbQuery)
	dbQuery.Count(&totalCount)
	totalPages := int(math.Ceil(float64(totalCount) / float64(config.PerPage)))

	//page number
	currentPage := utility.CurrentPage(c)

	dbQuery.
		Preload("BrowseNode").
		Order("rate desc, discount_percent desc, id desc").
		Limit(config.PerPage).
		Offset((currentPage - 1) * config.PerPage).
		Find(&category.Products)

	H := DefaultH(c)
	H["Category"] = &category
	H["Breadcrumbs"] = category.Breadcrumbs()
	H["MetaDescription"] = category.GetMetaDescription()
	H["MetaKeywords"] = category.GetMetaKeywords()
	H["Sidebar"] = true
	H["Title"] = category.Title
	H["Pagination"] = utility.Paginator(currentPage, totalPages, c.Request.URL)
	c.HTML(200, "categories/show", H)
}

//CategoriesGet processes GET /categories
func CategoriesGet(c *gin.Context) {
	var categories []models.Category
	models.DB.Find(&categories)

	for i := range categories {
		models.DB.Model(models.Product{}).Where("category_id = ?", categories[i].ID).Count(&categories[i].ProductCount)
	}

	H := DefaultH(c)
	H["Title"] = "Categories"
	H["Categories"] = categories
	c.HTML(200, "admin/categories/index", H)
}

//CategoriesNewGet processes new category request
func CategoriesNewGet(c *gin.Context) {
	var topLevelCategories []models.Category
	models.DB.Where("parent_id is NULL").Find(&topLevelCategories)

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	category := models.Category{}
	H := DefaultH(c)
	H["Title"] = "New category"
	H["TopLevelCategories"] = topLevelCategories
	H["Category"] = &category
	H["Flash"] = flashes
	c.HTML(200, "admin/categories/new", H)
}

//CategoriesNewPost processes create category request
func CategoriesNewPost(c *gin.Context) {
	category := models.Category{}
	if err := c.ShouldBind(&category); err != nil {
		sessionErrorAndRedirect(c, err, "/admin/new_category")
		return
	}

	if err := models.DB.Create(&category).Error; err != nil {
		sessionErrorAndRedirect(c, err, "/admin/new_category")
		return
	}
	if category.Submit == submitAndViewTitle() {
		c.Redirect(http.StatusSeeOther, category.GetURL())
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/categories")
}

//CategoriesEditGet processes edit category request
func CategoriesEditGet(c *gin.Context) {
	id := c.Param("id")

	category := models.Category{}
	models.DB.First(&category, id)

	if category.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	models.DB.Model(models.Product{}).Where("category_id = ?", category.ID).Count(&category.ProductCount)

	var topLevelCategories []models.Category
	models.DB.Where("parent_id is NULL").Find(&topLevelCategories)

	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	H := DefaultH(c)
	H["Title"] = category.Title
	H["TopLevelCategories"] = topLevelCategories
	H["Category"] = &category
	H["Flash"] = flashes
	c.HTML(200, "admin/categories/edit", H)
}

//CategoriesEditPost processes update category request
func CategoriesEditPost(c *gin.Context) {
	id := c.Param("id")

	category := models.Category{}
	models.DB.First(&category, id)

	if category.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	vm := models.Category{}
	if err := c.ShouldBind(&vm); err != nil {
		sessionErrorAndRedirect(c, err, fmt.Sprintf("/admin/edit_category/%d", category.ID))
		return
	}
	category.Title, category.Description, category.ParentID = vm.Title, vm.Description, vm.ParentID

	if err := models.DB.Save(&category).Error; err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/edit_category/%d", category.ID))
		return
	}
	if vm.Submit == submitAndViewTitle() {
		c.Redirect(http.StatusSeeOther, category.GetURL())
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/categories")
}

//CategoriesDeletePost processes delete category request
func CategoriesDeletePost(c *gin.Context) {
	id := c.Param("id")

	category := models.Category{}
	models.DB.First(&category, id)

	if category.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}
	if err := models.DB.Delete(&category).Error; err != nil {
		c.HTML(500, "errors/500", gin.H{"Error": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/categories")
}
