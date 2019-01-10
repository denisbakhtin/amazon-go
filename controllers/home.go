package controllers

import (
	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/gin-gonic/gin"
)

//HomeGet handles GET / route
func HomeGet(c *gin.Context) {
	var nodes []models.BrowseNode
	//top 7 tags
	models.DB.
		Where("product_count > 7 and parent_id is null").
		Order("product_count desc").
		Limit(7).
		Find(&nodes)

	//load products
	for i := range nodes {
		nodes[i].LoadAllChildren()
		ids := nodes[i].AppendIDs(models.SELFANDCHILDREN)
		models.DB.
			Preload("BrowseNode").
			Preload("Brand").
			Where("browse_node_id IN(?) and available=true", ids).
			Order("discount_percent desc, id desc").
			Limit(8).
			Find(&nodes[i].Products)
	}

	H := DefaultH(c)
	H["Title"] = config.SiteTitle
	H["Nodes"] = nodes
	H["HomeTopProducts"] = homeTopProducts()
	c.HTML(200, "home/show", H)
}
