package controllers

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/denisbakhtin/amazon-go/utility"
	"github.com/gin-gonic/gin"
)

//SearchGet processes GET /search
func SearchGet(c *gin.Context) {
	if len(c.Query("query")) == 0 {
		c.Redirect(303, "/")
		return
	}

	var nodeIDs []uint64
	node := models.BrowseNode{}

	if nodeID := c.Query("category_id"); nodeID != "" && nodeID != "0" {
		models.DB.First(&node, nodeID)
		node.LoadAllChildren()
		nodeIDs = node.AppendIDs(models.SELFANDCHILDREN)
	}

	query := c.Query("query")
	//prepare search string and split into array of words
	query = strings.ToLower(query)
	re := regexp.MustCompile("\\s+")
	query = re.ReplaceAllString(query, " ")
	re = regexp.MustCompile("[^\\w\\s-]")
	query = re.ReplaceAllString(query, "")
	searchTerms := strings.Split(query, " ")

	//page number
	currentPage := 1
	if pageStr := c.Query("page"); pageStr != "" {
		currentPage, _ = strconv.Atoi(pageStr)
	}
	currentPage = int(math.Max(float64(1), float64(currentPage)))

	//postgresql full text search
	dbQuery := models.DB.Where("available = true AND to_tsvector('english', title) @@ to_tsquery('english', ?)", strings.Join(searchTerms, " & "))

	if len(nodeIDs) > 0 {
		dbQuery = dbQuery.Where("browse_node_id IN(?)", nodeIDs)
	}

	totalCount := 0
	dbQuery.Model(models.Product{}).Count(&totalCount)
	totalPages := int(math.Ceil(float64(totalCount) / float64(config.PerPage)))

	var products []models.Product
	dbQuery.
		Preload("BrowseNode").
		Order("rate desc, discount_percent desc, id desc").
		Limit(config.PerPage).
		Offset((currentPage - 1) * config.PerPage).
		Find(&products) //order by rating + apply pagination

	H := DefaultH(c)
	H["Products"] = products
	H["Pagination"] = utility.Paginator(currentPage, totalPages, c.Request.URL)
	H["SearchString"] = strings.Join(searchTerms, " ")
	H["Title"] = fmt.Sprintf("%s - search results", strings.Join(searchTerms, " "))
	c.HTML(200, "search/show", H)
}
