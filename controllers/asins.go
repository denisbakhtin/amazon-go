package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/denisbakhtin/amazon-go/aws"
	"github.com/denisbakhtin/amazon-go/models"
)

//ProcessedAsinsGet shows all processed asins
func ProcessedAsinsGet(c *gin.Context) {
	var asins []models.ProcessedAsin
	models.DB.Order("id asc").Find(&asins)

	H := DefaultH(c)
	H["Title"] = "Processed Asins"
	H["Asins"] = asins
	c.HTML(200, "admin/asins/processed_index", H)
}

//QueuedAsinsGet shows all queued asins
func QueuedAsinsGet(c *gin.Context) {
	var asins []models.QueuedAsin
	models.DB.Order("id asc").Preload("Feed").Find(&asins)

	H := DefaultH(c)
	H["Title"] = "Queued Asins"
	H["Asins"] = asins
	c.HTML(200, "admin/asins/queued_index", H)
}

//QueuedAsinsNewGet processes new queued_asin request
func QueuedAsinsNewGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	asin := models.QueuedAsin{}
	H := DefaultH(c)
	H["Title"] = "Queue Asin"
	H["Asin"] = &asin
	H["Flash"] = flashes
	c.HTML(200, "admin/asins/queued_new", H)
}

//QueuedAsinsNewPost processes create queued_asin request
func QueuedAsinsNewPost(c *gin.Context) {
	asin := models.QueuedAsin{}
	if err := c.ShouldBind(&asin); err != nil {
		sessionErrorAndRedirect(c, err, "/admin/new_queued_asin")
		return
	}
	if err := models.DB.Create(&asin).Error; err != nil {
		sessionErrorAndRedirect(c, err, "/admin/new_queued_asin")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/queued_asins")
}

//QueuedAsinsNewProductGet processes new queued_product_id request
func QueuedAsinsNewProductGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	H := DefaultH(c)
	H["Title"] = "Queue Product ID"
	H["Flash"] = flashes
	c.HTML(200, "admin/asins/queued_new_product_id", H)
}

//QueuedAsinsNewProductPost processes create queued_product_id request
func QueuedAsinsNewProductPost(c *gin.Context) {
	pID := models.QueuedProductID{}
	if err := c.ShouldBind(&pID); err != nil {
		sessionErrorAndRedirect(c, err, "/admin/new_queued_product_id")
		return
	}
	product := models.Product{}
	models.DB.First(&product, pID.ProductID)
	if product.ID > 0 {
		asin := models.QueuedAsin{Asin: product.Asin}
		if err := models.DB.Create(&asin).Error; err != nil {
			sessionErrorAndRedirect(c, err, "/admin/new_queued_product_id")
			return
		}
	}
	c.Redirect(http.StatusSeeOther, "/admin/queued_asins")
}

//QueuedAsinsDeletePost processes delete queued_asin request
func QueuedAsinsDeletePost(c *gin.Context) {
	id := c.Param("id")

	asin := models.QueuedAsin{}
	models.DB.First(&asin, id)

	if asin.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}
	if err := models.DB.Delete(&asin).Error; err != nil {
		c.HTML(500, "errors/500", gin.H{"Error": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/queued_asins")
}

//ProcessedAsinsDeletePost processes delete processed_asin request
func ProcessedAsinsDeletePost(c *gin.Context) {
	id := c.Param("id")

	asin := models.ProcessedAsin{}
	models.DB.First(&asin, id)

	if asin.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}
	if err := models.DB.Delete(&asin).Error; err != nil {
		c.HTML(500, "errors/500", gin.H{"Error": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/processed_asins")
}

//QueuedAsinsClearPost clears asin queue
func QueuedAsinsClearPost(c *gin.Context) {
	aws.ClearQueuedAsins()
	c.Redirect(http.StatusSeeOther, "/admin/queued_asins")
}

//ProcessedAsinsClearPost clears processed asins table
func ProcessedAsinsClearPost(c *gin.Context) {
	aws.ClearProcessedAsins()
	if len(c.Request.Referer()) > 0 {
		//if called from queued asins page, return back there
		c.Redirect(http.StatusSeeOther, c.Request.Referer())
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/processed_asins")
}

//QueueAsinsPost queues available products
func QueueAsinsPost(c *gin.Context) {
	aws.QueueAvailableAsins()
	c.Redirect(http.StatusSeeOther, "/admin/queued_asins")
}

//QueueAllAsinsPost queues available & unavailable products
func QueueAllAsinsPost(c *gin.Context) {
	aws.QueueAvailableAsins()
	aws.QueueUnavailableAsins()
	c.Redirect(http.StatusSeeOther, "/admin/queued_asins")
}
