package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/denisbakhtin/amazon-go/utility"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
	"gopkg.in/go-playground/validator.v8"
)

const (
	accountIDKey   = "account_id"
	accountRoleKey = "account_role"
)

var router *gin.Engine

//InitRoutes initializes web site routes
func InitRoutes() {

	gin.SetMode(config.GetMode())
	router = gin.Default()
	router.SetHTMLTemplate(loadTemplates())

	router.Use(static.Serve("/", static.LocalFile(config.PublicPath, false)))
	router.Use(sessions.Sessions("aws-session", config.SessionStore))
	router.Use(contextData())

	//cart requests without CSRF protection
	router.GET("/cart", CartGet)
	router.GET("/cart/get", CartJSONGet)
	router.POST("/cart/add/:id", CartAdd)
	//router.POST("/cart/add", CartAddJSONPost)
	router.POST("/cart/update", CartUpdatePost)
	router.GET("/cart/checkout", CartCheckoutGet)

	//setup csrf protection
	router.Use(csrf.Middleware(csrf.Options{
		Secret: config.CsrfSecret,
		ErrorFunc: func(c *gin.Context) {
			log.Println("CSRF token mismatch")
			ShowErrorPage(c, 400, fmt.Errorf("CSRF token mismatch"))
			c.Abort()
		},
	}))

	router.NoRoute(Error404)

	router.GET("/", HomeGet)
	router.GET("/products/:id/:slug", ProductGet)
	//router.GET("/product_reviews/:id", ProductReviewsGet)
	router.GET("/variations/:asin", VariationJSONGet)
	router.GET("/tags/:id/:slug", TagGet)
	router.GET("/brands/:id/:slug", BrandGet)
	router.GET("/b/:id/:slug", BindingGet)
	router.GET("/pages/:id/:slug", PageGet)
	router.GET("/search", SearchGet)
	router.GET("/signin", SignInGet)
	router.POST("/signin", SignInPost)
	if config.SignUpEnabled {
		router.GET("/signup", SignUpGet)
		router.POST("/signup", SignUpPost)
	}
	router.Any("/signout", SignOutPost)

	admin := router.Group("/admin", authenticatedAdmin())
	{
		admin.GET("/", DashboardGet)
		admin.GET("/categories", CategoriesGet)
		admin.GET("/new_category", CategoriesNewGet)
		admin.POST("/new_category", CategoriesNewPost)
		admin.GET("/edit_category/:id", CategoriesEditGet)
		admin.POST("/edit_category/:id", CategoriesEditPost)
		admin.POST("/delete_category/:id", CategoriesDeletePost)

		admin.GET("/pages", PagesGet)
		admin.GET("/new_page", PagesNewGet)
		admin.POST("/new_page", PagesNewPost)
		admin.GET("/edit_page/:id", PagesEditGet)
		admin.POST("/edit_page/:id", PagesEditPost)
		admin.POST("/delete_page/:id", PagesDeletePost)

		admin.GET("/accounts", AccountsGet)
		admin.GET("/new_account", AccountsNewGet)
		admin.POST("/new_account", AccountsNewPost)
		admin.GET("/edit_account/:id", AccountsEditGet)
		admin.POST("/edit_account/:id", AccountsEditPost)
		admin.POST("/delete_account/:id", AccountsDeletePost)

		admin.GET("/processed_specifications", ProcessedSpecificationsGet)
		admin.GET("/queued_specifications", QueuedSpecificationsGet)
		admin.GET("/new_queued_specification", QueuedSpecificationsNewGet)
		admin.POST("/new_queued_specification", QueuedSpecificationsNewPost)
		admin.POST("/delete_queued_specification/:id", QueuedSpecificationsDeletePost)
		admin.POST("/clear_queued_specifications", QueuedSpecificationsClearPost)
		admin.POST("/delete_processed_specification/:id", ProcessedSpecificationsDeletePost)
		admin.POST("/clear_processed_specifications", ProcessedSpecificationsClearPost)
		admin.POST("/queue_specifications", QueueSpecificationsPost)

		admin.GET("/processed_asins", ProcessedAsinsGet)
		admin.GET("/queued_asins", QueuedAsinsGet)
		admin.GET("/new_queued_asin", QueuedAsinsNewGet)
		admin.POST("/new_queued_asin", QueuedAsinsNewPost)
		admin.GET("/new_queued_product_id", QueuedAsinsNewProductGet)
		admin.POST("/new_queued_product_id", QueuedAsinsNewProductPost)
		admin.POST("/delete_queued_asin/:id", QueuedAsinsDeletePost)
		admin.POST("/clear_queued_asins", QueuedAsinsClearPost)
		admin.POST("/delete_processed_asin/:id", ProcessedAsinsDeletePost)
		admin.POST("/clear_processed_asins", ProcessedAsinsClearPost)
		admin.POST("/queue_asins", QueueAsinsPost)
		admin.POST("/queue_all_asins", QueueAllAsinsPost)

		admin.POST("/upload", UploadPost)

		admin.GET("/browse_nodes", BrowseNodesGet)
		admin.GET("/browse_node/:id", BrowseNodeGet)
	}
}

//authenticated is authentication middleware, enabled by router for protected routes
func authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		account, _ := c.Get("account")
		if account == nil || account.(models.Account).ID == 0 {
			c.AbortWithStatus(403)
		}
		c.Next()
	}
}

//authenticatedAdmin is a middleware, restricting access to all but admins
func authenticatedAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		account, _ := c.Get("account")
		if account == nil || account.(models.Account).ID == 0 || account.(models.Account).Role != config.AdminRole {
			c.AbortWithStatus(403)
		}
		c.Next()
	}
}

//contextData is a middleware, setting shared context data
func contextData() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get(accountIDKey) != nil {
			account := models.Account{}
			models.DB.First(&account, session.Get(accountIDKey))
			if account.ID != 0 {
				c.Set("account", account)
			}
		}
	}
}

//GetRouter returns gin router
func GetRouter() *gin.Engine {
	return router
}

const requiredFieldErrorMsg = "'%s' is required."
const fieldErrorMsg = "Field validation for '%s' failed, because it is %s."

//sessionErrorAndRedirect adds error to session flashes
func sessionErrorAndRedirect(c *gin.Context, err error, url string) {
	session := sessions.Default(c)
	switch ert := err.(type) {
	case validator.ValidationErrors:
		//Replace default gorm validator error strings
		for _, er := range ert {
			if er.Tag == "required" {
				session.AddFlash(fmt.Sprintf(requiredFieldErrorMsg, utility.SplitCamelWords(er.Field)))
			} else {
				session.AddFlash(fmt.Sprintf(fieldErrorMsg, utility.SplitCamelWords(er.Field), er.Tag))
			}
		}
	default:
		session.AddFlash(err.Error())
	}
	session.Save()
	c.Redirect(http.StatusSeeOther, url)
}
