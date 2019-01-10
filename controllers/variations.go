package controllers

import (
	"bytes"
	"log"

	"github.com/denisbakhtin/amazon-go/models"
	"github.com/gin-gonic/gin"
)

//VariationJSONGet processes get /variations/:asin json request
func VariationJSONGet(c *gin.Context) {
	asin := c.Param("asin")

	variation := models.Variation{}
	models.DB.Where("asin = ?", asin).First(&variation)
	if variation.ID == 0 {
		c.JSON(404, nil)
	}
	var images, description, offer bytes.Buffer
	if err := tmpl.ExecuteTemplate(&images, "products/product_images", variation); err != nil {
		log.Println(err)
		c.JSON(500, nil)
		return
	}
	if err := tmpl.ExecuteTemplate(&description, "products/short_description", variation); err != nil {
		log.Println(err)
		c.JSON(500, nil)
		return
	}
	if err := tmpl.ExecuteTemplate(&offer, "products/offer_details", variation); err != nil {
		log.Println(err)
		c.JSON(500, nil)
		return
	}
	c.JSON(200, gin.H{
		"Images":      images.String(),
		"Description": description.String(),
		"Offer":       offer.String(),
	})
}
