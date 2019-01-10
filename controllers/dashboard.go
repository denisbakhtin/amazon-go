package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//DashboardGet handles GET /admin route
func DashboardGet(c *gin.Context) {
	/* H := DefaultH(c)
	H["Title"] = "Admin dashboard"
	c.HTML(200, "admin/dashboard/show", H) */
	c.Redirect(http.StatusSeeOther, "/admin/categories")
}
