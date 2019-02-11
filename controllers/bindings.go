package controllers

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/denisbakhtin/amazon-go/utility"
)

//BindingGet processes GET /b/%d/%s route
func BindingGet(c *gin.Context) {
	id := c.Param("id")

	binding := models.Binding{}
	models.DB.First(&binding, id)

	if binding.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	//redirect to canonical url
	if c.Request.URL.Path != binding.GetURL() {
		c.Redirect(http.StatusMovedPermanently, binding.GetURL())
		return
	}

	totalCount := 0
	dbQuery := models.DB.Model(models.Product{}).Where("binding_id = ? AND available = true", id)
	dbQuery = applyProductFilters(c, dbQuery)
	dbQuery.Count(&totalCount)
	totalPages := int(math.Ceil(float64(totalCount) / float64(config.PerPage)))

	//page number
	currentPage := utility.CurrentPage(c)

	var products []models.Product
	dbQuery.
		Preload("BrowseNode").
		Order("rate desc, discount_percent desc, id desc").
		Limit(config.PerPage).
		Offset((currentPage - 1) * config.PerPage).
		Find(&products)

	H := DefaultH(c)
	H["Title"] = fmt.Sprintf("%s Category", binding.Title)
	H["Binding"] = &binding
	H["Products"] = products
	H["BodyClass"] = "mw-1400"
	H["Sidebar"] = true
	H["MetaKeywords"] = binding.GetMetaKeywords()
	H["MetaDescription"] = binding.GetMetaDescription()
	H["Pagination"] = utility.Paginator(currentPage, totalPages, c.Request.URL)
	c.HTML(200, "bindings/show", H)
}
