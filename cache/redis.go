package cache

import (
	"fmt"
	"log"
	"time"

	"encoding/json"

	"github.com/denisbakhtin/amazon-go/config"
	"github.com/denisbakhtin/amazon-go/models"
	"github.com/go-redis/redis"
)

//Client is a redis client
var Client *redis.Client

const (
	homeNodesKey         = "home_nodes"
	homeNodeChildrenFKey = "home_node_children_%d"
	homeProductsKey      = "home_products"
	nodeBrandsFKey       = "node_brands_%d"
	brandNodesFKey       = "brand_nodes_%d"
)

//Init explicitly initializes redis cache store
func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.RedisConnectionString,
		Password: config.RedisPassword, // no password set
		DB:       0,                    // use default DB
	})

	if _, err := Client.Ping().Result(); err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println("Redis cache has been initialized")
}

//UpdateAll10Min updates all cache strings with 10 minutes lifetime
func UpdateAll10Min() {
	SetHomeNodes()
	SetHomeTopProducts()
}

//UpdateAll30Min updates all cache strings with 30 minutes lifetime
func UpdateAll30Min() {
	SetNodesBrands(models.TopNodes())
	SetBrandsNodes()
}

//GetHomeNodes reads browse nodes shown at home page from cache, if absent, creates cache
func GetHomeNodes() []models.BrowseNode {
	s, err := Client.Get(homeNodesKey).Result()
	if err != nil {
		return SetHomeNodes()
	}
	var result []models.BrowseNode
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		log.Println("Error unmarshalling home nodes cache: ", err)
		return SetHomeNodes()
	}
	return result
}

//SetHomeNodes caches browse nodes shown at home page
//10 minutes lifetime
func SetHomeNodes() []models.BrowseNode {
	var nodes []models.BrowseNode
	models.DB.
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
		Where("product_count > 7 and parent_id is null").
		Order("product_count desc").
		Limit(7).
		Find(&nodes)

	for i := range nodes {
		nodes[i].AllChildrenLoaded = true
		nodes[i].ThreeChildren = GetThreeChildren(&nodes[i])
		ids := nodes[i].AppendIDs(models.SELFANDCHILDREN)
		//clear children to minimize cache size
		nodes[i].Children = nil
		models.DB.
			Select("id,title,image,discount_percent,regular_price,special_price,browse_node_id,brand_id").
			Preload("BrowseNode").
			Preload("Brand").
			Where("browse_node_id IN(?) and available=true", ids).
			Order("discount_percent desc, id desc").
			Limit(8).
			Find(&nodes[i].Products)
	}
	buf, err := json.Marshal(nodes)
	if err != nil {
		log.Println("Error marshalling home nodes for cache: ", err)
		return nodes
	}
	if err := Client.Set(homeNodesKey, string(buf), time.Duration(10)*time.Minute).Err(); err != nil {
		log.Println("Error storing home nodes json into cache: ", err)
	}
	return nodes
}

//GetThreeChildren reads 3 child browse nodes from cache, if absent, creates cache
func GetThreeChildren(b *models.BrowseNode) []models.BrowseNode {
	s, err := Client.Get(fmt.Sprintf(homeNodeChildrenFKey, b.ID)).Result()
	if err != nil {
		return SetThreeChildren(b)
	}
	var result []models.BrowseNode
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		log.Println("Error unmarshalling child nodes from cache: ", err)
		return SetThreeChildren(b)
	}
	return result
}

//SetThreeChildren caches child browse nodes shown at home page
//10 minutes lifetime
func SetThreeChildren(b *models.BrowseNode) []models.BrowseNode {
	if !b.AllChildrenLoaded {
		b.LoadAllChildren()
	}
	nodes := make([]models.BrowseNode, 0, 3)
	nodes = appendChildren(3, b, nodes)
	//clear loaded children now
	for i := range nodes {
		nodes[i].Children = nil
	}

	buf, err := json.Marshal(nodes)
	if err != nil {
		log.Println("Error marshalling child nodes for cache: ", err)
		return nodes
	}
	if err := Client.Set(fmt.Sprintf(homeNodeChildrenFKey, b.ID), string(buf), time.Duration(10)*time.Minute).Err(); err != nil {
		log.Println("Error storing child nodes json into cache: ", err)
	}
	return nodes
}

func appendChildren(number int, b *models.BrowseNode, result []models.BrowseNode) []models.BrowseNode {
	//append direct children first
	for i := 0; len(result) < number && i < len(b.Children); i++ {
		if b.Children[i].OwnProductCount > 0 {
			result = append(result, b.Children[i])
		}
	}
	if len(result) < number {
		for i := range b.Children {
			result = appendChildren(number, &b.Children[i], result)
		}
	}
	return result
}

//GetHomeTopProducts reads top products for home page from cache, if absent, creates cache
func GetHomeTopProducts() []models.Product {
	s, err := Client.Get(homeProductsKey).Result()
	if err != nil {
		return SetHomeTopProducts()
	}
	var result []models.Product
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		log.Println("Error unmarshalling home top products cache: ", err)
		return SetHomeTopProducts()
	}
	return result
}

//SetHomeTopProducts caches top products shown at home page
//10 minutes lifetime
func SetHomeTopProducts() []models.Product {
	var ids []uint64
	models.DB.Table("products").Where("discount_percent > 0 and available = true and image != ?", "").Select("min(id) as idd").Group("browse_node_id").Order("browse_node_id desc").Limit(30).Pluck("idd", &ids)
	var products []models.Product
	models.DB.Where("id IN(?)", ids).Order("id asc").Find(&products)

	buf, err := json.Marshal(products)
	if err != nil {
		log.Println("Error marshalling home top products for cache: ", err)
		return products
	}
	if err := Client.Set(homeProductsKey, string(buf), time.Duration(10)*time.Minute).Err(); err != nil {
		log.Println("Error storing home top products json into cache: ", err)
	}
	return products
}

//GetNodeBrands reads node's brands from cache, if absent, creates cache
func GetNodeBrands(node *models.BrowseNode) []models.Brand {
	s, err := Client.Get(fmt.Sprintf(nodeBrandsFKey, node.ID)).Result()
	if err != nil {
		return SetNodeBrands(node)
	}
	var result []models.Brand
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		log.Println("Error unmarshalling node's brands from cache: ", err)
		return SetNodeBrands(node)
	}
	return result
}

//SetNodeBrands caches node's brands
//30 minutes lifetime
func SetNodeBrands(node *models.BrowseNode) []models.Brand {
	node.LoadAllChildren()
	ids := node.AppendIDs(models.SELFANDCHILDREN)
	var brandIDs []uint64
	models.DB.Table("products").Where("available = true AND browse_node_id IN(?)", ids).Select("brand_id").Group("brand_id").Pluck("brand_id", &brandIDs)
	var brands []models.Brand
	models.DB.Where("id IN(?)", brandIDs).Order("title asc").Find(&brands)

	buf, err := json.Marshal(brands)
	if err != nil {
		log.Println("Error marshalling node's brands for cache: ", err)
		return brands
	}
	if err := Client.Set(fmt.Sprintf(nodeBrandsFKey, node.ID), string(buf), time.Duration(30)*time.Minute).Err(); err != nil {
		log.Println("Error storing node's brands json into cache: ", err)
	}
	return brands
}

//SetNodesBrands caches each node's brands
//All children have been preloaded, just set the flag
func SetNodesBrands(nodes []models.BrowseNode) {
	for i := range nodes {
		nodes[i].AllChildrenLoaded = true //flag
		SetNodeBrands(&nodes[i])
		SetNodesBrands(nodes[i].Children)
	}
}

//GetBrandNodes reads brands's nodes from cache, if absent, creates cache
func GetBrandNodes(brand *models.Brand) []models.BrowseNode {
	s, err := Client.Get(fmt.Sprintf(brandNodesFKey, brand.ID)).Result()
	if err != nil {
		return SetBrandNodes(brand)
	}
	var result []models.BrowseNode
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		log.Println("Error unmarshalling brands's nodes from cache: ", err)
		return SetBrandNodes(brand)
	}
	return result
}

//SetBrandNodes caches brands's nodes
//30 minutes lifetime
func SetBrandNodes(brand *models.Brand) []models.BrowseNode {
	var nodeIDs []uint64
	models.DB.Table("products").Where("available = true AND brand_id = ?", brand.ID).Select("browse_node_id").Group("browse_node_id").Pluck("browse_node_id", &nodeIDs)
	var nodes []models.BrowseNode
	models.DB.Where("id IN(?)", nodeIDs).Order("title asc").Find(&nodes)

	buf, err := json.Marshal(nodes)
	if err != nil {
		log.Println("Error marshalling brand's nodes for cache: ", err)
		return nodes
	}
	if err := Client.Set(fmt.Sprintf(brandNodesFKey, brand.ID), string(buf), time.Duration(30)*time.Minute).Err(); err != nil {
		log.Println("Error storing brand's nodes json into cache: ", err)
	}
	return nodes
}

//SetBrandsNodes caches each brand's nodes
func SetBrandsNodes() {
	var brands []models.Brand
	models.DB.Find(&brands)
	for i := range brands {
		SetBrandNodes(&brands[i])
	}
}
