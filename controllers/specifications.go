package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/denisbakhtin/amazon-go/aws"
	"github.com/denisbakhtin/amazon-go/models"
)

//ProcessedSpecificationsGet shows all processed specifications
func ProcessedSpecificationsGet(c *gin.Context) {
	var specifications []models.ProcessedSpecification
	models.DB.Order("id asc").Find(&specifications)

	H := DefaultH(c)
	H["Title"] = "Processed Specifications"
	H["Specifications"] = specifications
	c.HTML(200, "admin/specifications/processed_index", H)
}

//QueuedSpecificationsGet shows all queued specifications
func QueuedSpecificationsGet(c *gin.Context) {
	var specifications []models.QueuedSpecification
	models.DB.Order("id asc").Preload("Product").Find(&specifications)

	H := DefaultH(c)
	H["Title"] = "Queued Specifications"
	H["Specifications"] = specifications
	c.HTML(200, "admin/specifications/queued_index", H)
}

//QueuedSpecificationsNewGet processes new queued_specification request
func QueuedSpecificationsNewGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	specification := models.QueuedSpecification{}
	H := DefaultH(c)
	H["Title"] = "New Queued Specification"
	H["Specification"] = &specification
	H["Flash"] = flashes
	c.HTML(200, "admin/specifications/queued_new", H)
}

//QueuedSpecificationsNewPost processes create queued_specification request
func QueuedSpecificationsNewPost(c *gin.Context) {
	specification := models.QueuedSpecification{}
	if err := c.ShouldBind(&specification); err != nil {
		sessionErrorAndRedirect(c, err, "/admin/new_queued_specification")
		return
	}
	if err := models.DB.Create(&specification).Error; err != nil {
		sessionErrorAndRedirect(c, err, "/admin/new_queued_specification")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/queued_specifications")
}

//QueuedSpecificationsDeletePost processes delete queued_specification request
func QueuedSpecificationsDeletePost(c *gin.Context) {
	id := c.Param("id")

	specification := models.QueuedSpecification{}
	models.DB.First(&specification, id)

	if specification.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}
	if err := models.DB.Delete(&specification).Error; err != nil {
		c.HTML(500, "errors/500", gin.H{"Error": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/queued_specifications")
}

//ProcessedSpecificationsDeletePost processes delete processed_specification request
func ProcessedSpecificationsDeletePost(c *gin.Context) {
	id := c.Param("id")

	specification := models.ProcessedSpecification{}
	models.DB.First(&specification, id)

	if specification.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}
	if err := models.DB.Delete(&specification).Error; err != nil {
		c.HTML(500, "errors/500", gin.H{"Error": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/processed_specifications")
}

//QueuedSpecificationsClearPost clears specification queue
func QueuedSpecificationsClearPost(c *gin.Context) {
	aws.ClearQueuedSpecifications()
	c.Redirect(http.StatusSeeOther, "/admin/queued_specifications")
}

//ProcessedSpecificationsClearPost clears processed specifications table
func ProcessedSpecificationsClearPost(c *gin.Context) {
	aws.ClearProcessedSpecifications()
	c.Redirect(http.StatusSeeOther, "/admin/processed_specifications")
}

//QueueSpecificationsPost queues available products
func QueueSpecificationsPost(c *gin.Context) {
	aws.QueueAvailableSpeficications()
	c.Redirect(http.StatusSeeOther, "/admin/queued_specifications")
}
