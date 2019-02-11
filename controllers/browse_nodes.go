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

//BrowseNodeGet processes GET /admin/browse_node/1
func BrowseNodeGet(c *gin.Context) {
	id := c.Param("id")

	node := models.BrowseNode{}
	models.DB.Preload("Products").First(&node, id)

	if node.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	H := DefaultH(c)
	H["Title"] = fmt.Sprintf("#%s", node.Title)
	H["BrowseNode"] = &node
	c.HTML(200, "admin/browse_nodes/show", H)
}

//TagGet processes GET /tags/%d/%s route
func TagGet(c *gin.Context) {
	id := c.Param("id")

	node := models.BrowseNode{}
	models.DB.First(&node, id)

	if node.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	//redirect to canonical url
	if c.Request.URL.Path != node.GetURL() {
		c.Redirect(http.StatusMovedPermanently, node.GetURL())
		return
	}

	node.LoadAllChildren()
	node.LoadAllParents()
	ids := node.AppendIDs(models.SELFANDCHILDREN)
	totalCount := 0
	dbQuery := models.DB.Model(models.Product{}).Where("browse_node_id IN(?) AND available = true", ids)
	dbQuery = applyProductFilters(c, dbQuery)
	dbQuery.Count(&totalCount)
	totalPages := int(math.Ceil(float64(totalCount) / float64(config.PerPage)))

	//page number
	currentPage := utility.CurrentPage(c)

	dbQuery.
		Preload("BrowseNode").
		Preload("Brand").
		Order("rate desc, discount_percent desc, id desc").
		Limit(config.PerPage).
		Offset((currentPage - 1) * config.PerPage).
		Find(&node.Products)

	if (len(node.Products) < config.PerPage) && (currentPage == 1) {
		topParent := node.TopParent()
		topParent.LoadAllChildren()
		allids := topParent.AppendIDs(models.SELFANDCHILDREN)
		similarids := utility.SubtractUint64Array(allids, ids)
		models.DB.
			Preload("BrowseNode").
			Where("browse_node_id IN(?) AND available = true", similarids).
			Order("rate desc, discount_percent desc, id desc").
			Limit(config.PerPage).
			Find(&node.SimilarProducts)
	}

	H := DefaultH(c)
	H["Title"] = node.Title
	H["Tag"] = &node
	H["BodyClass"] = "mw-1400"
	H["Sidebar"] = true
	H["MetaKeywords"] = node.GetMetaKeywords()
	H["MetaDescription"] = node.GetMetaDescription()
	H["SidebarBrands"] = cache.GetNodeBrands(&node)
	H["Pagination"] = utility.Paginator(currentPage, totalPages, c.Request.URL)
	c.HTML(200, "tags/show", H)
}

//BrowseNodesGet shows all browse_nodes
func BrowseNodesGet(c *gin.Context) {
	var nodes []models.BrowseNode
	if c.Query("only_top") == "1" {
		models.DB.Order("parent_id asc, id asc").
			Where("parent_id is null").Order("id asc").Find(&nodes)
	} else {
		models.DB.Order("parent_id asc, id asc").
			Preload("Children").
			Preload("Children.Children").
			Preload("Children.Children.Children").
			Preload("Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Preload("Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children.Children").
			Where("parent_id is null").Order("id asc").Find(&nodes)
	}

	H := DefaultH(c)
	H["Title"] = "Browse nodes"
	H["BrowseNodes"] = nodes
	H["OnlyTop"] = (c.Query("only_top") == "1")
	c.HTML(200, "admin/browse_nodes/index", H)
}
