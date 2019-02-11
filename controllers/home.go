package controllers

import (
	"fmt"

	"github.com/denisbakhtin/amazon-go/cache"
	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/gin-gonic/gin"
)

//HomeGet handles GET / route
func HomeGet(c *gin.Context) {
	H := DefaultH(c)
	H["Title"] = fmt.Sprintf("Shop special offers and deals on %s", config.SiteTitle)
	H["MetaKeywords"] = "Special offers, International shipping deals, Discounts, Discount centre"
	H["MetaDescription"] = "Buy popular products via gateway shopping centre with international shipping. Low prices on brand items!"
	H["Nodes"] = cache.GetHomeNodes()
	H["MenuNodes"] = models.MenuNodes()
	H["HomeTopProducts"] = cache.GetHomeTopProducts()
	c.HTML(200, "home/show", H)
}
