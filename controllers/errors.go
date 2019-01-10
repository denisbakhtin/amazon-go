package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Error404 shows custom 404 error page
func Error404(c *gin.Context) {
	ShowErrorPage(c, http.StatusNotFound, nil)
}

//ShowErrorPage executes error template given its code
func ShowErrorPage(c *gin.Context, code int, err error) {
	H := DefaultH(c)
	H["Error"] = err
	H["Controller"] = "error"
	c.HTML(code, fmt.Sprintf("errors/%d", code), H)
}
