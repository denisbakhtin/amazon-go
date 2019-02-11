package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/denisbakhtin/amazon-go/utility"

	"github.com/denisbakhtin/amazon-go/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ProductGet processes GET /products/1/slug
func ProductGet(c *gin.Context) {
	id := c.Param("id")

	product := models.Product{}
	models.DB.
		Preload("Variations", func(db *gorm.DB) *gorm.DB {
			return db.Preload("ItemAttributes").Where("available = true").Order("variations.dim1_value ASC, variations.discount_percent DESC")
		}).
		Preload("Brand").
		Preload("Binding").
		Preload("BrowseNode").
		Preload("BrowseNode.Parent").
		Preload("BrowseNode.Parent.Parent").
		Preload("BrowseNode.Parent.Parent.Parent").
		Preload("BrowseNode.Parent.Parent.Parent.Parent").
		Preload("BrowseNode.Parent.Parent.Parent.Parent.Parent").
		Preload("BrowseNode.Parent.Parent.Parent.Parent.Parent.Parent").
		Preload("BrowseNode.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
		Preload("BrowseNode.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
		Preload("BrowseNode.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
		Preload("BrowseNode.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
		Preload("BrowseNode.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent.Parent").
		First(&product, id)

	if product.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	//redirect to canonical url
	if c.Request.RequestURI != product.GetURL() {
		c.Redirect(http.StatusMovedPermanently, product.GetURL())
		return //return is needed to prevent further execution
	}

	H := DefaultH(c)
	H["Product"] = &product
	H["Title"] = product.Title
	H["MetaKeywords"] = product.GetMetaKeywords()
	H["MetaDescription"] = product.GetMetaDescription()
	H["SimilarProducts"] = product.Similar()
	c.HTML(200, "products/show", H)
}

//ProductReviewsGet processes GET /product_reviews/1
//obsolete
func ProductReviewsGet(c *gin.Context) {
	id := c.Param("id")

	product := models.Product{}
	models.DB.Select("customer_reviews").First(&product, id)

	if len(product.CustomerReviews) == 0 {
		c.String(201, "")
		return
	}

	//resp, err := http.Get(product.CustomerReviews)
	//emulate firefox browser :)
	resp, err := utility.GetAmazonURL(product.CustomerReviews)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		fmt.Println(err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		fmt.Println(err)
	}
	sanitizer := utility.ReviewsSanitizer()
	c.String(200, sanitizer.Sanitize(string(contents)))
}
