package controllers

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/denisbakhtin/amazon-go/cache"
	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/denisbakhtin/amazon-go/utility"
)

//BrandGet processes GET /brands/%d/%s route
func BrandGet(c *gin.Context) {
	id := c.Param("id")

	brand := models.Brand{}
	models.DB.First(&brand, id)

	if brand.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	//redirect to canonical url
	if c.Request.URL.Path != brand.GetURL() {
		c.Redirect(http.StatusMovedPermanently, brand.GetURL())
		return
	}

	totalCount := 0
	dbQuery := models.DB.Model(models.Product{}).Where("brand_id = ? AND available = true", id)
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
	H["Title"] = fmt.Sprintf("%s Brand", brand.Title)
	H["Brand"] = &brand
	H["Products"] = products
	H["BodyClass"] = "mw-1400"
	H["Sidebar"] = true
	H["SidebarTags"] = cache.GetBrandNodes(&brand)
	H["MetaKeywords"] = brand.GetMetaKeywords()
	H["MetaDescription"] = brand.GetMetaDescription()
	H["Pagination"] = utility.Paginator(currentPage, totalPages, c.Request.URL)
	c.HTML(200, "brands/show", H)
}
